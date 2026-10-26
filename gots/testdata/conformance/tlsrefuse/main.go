// tlsrefuse: what a TLS client or server must refuse it refuses, and what
// it must take it takes: certificates from an unknown authority, for
// another name, expired, or failing a check of the program's own, a client
// certificate a server cannot verify, and versions outside those allowed.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"time"
)

type issued struct {
	cert tls.Certificate
	leaf *x509.Certificate
}

var serial int64

// issue makes a certificate for name valid from nb to na, signed by ca, or
// by itself as a CA when ca is nil.
func issue(name string, ca *issued, nb, na time.Time, usage x509.ExtKeyUsage) *issued {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	serial++
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: name},
		NotBefore: nb, NotAfter: na, DNSNames: []string{name},
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{usage},
	}
	parent, signer := tmpl, any(key)
	if ca == nil {
		tmpl.IsCA, tmpl.BasicConstraintsValid = true, true
		tmpl.KeyUsage |= x509.KeyUsageCertSign
		tmpl.ExtKeyUsage = nil
	} else {
		parent, signer = ca.leaf, ca.cert.PrivateKey
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &key.PublicKey, signer)
	if err != nil {
		panic(err)
	}
	leaf, _ := x509.ParseCertificate(der)
	return &issued{tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, leaf}
}

func pool(cs ...*issued) *x509.CertPool {
	p := x509.NewCertPool()
	for _, c := range cs {
		p.AddCert(c.leaf)
	}
	return p
}

// try runs one handshake of a server with sc and a client with cc, and
// tells how each side ended: the client writing, the server echoing.
func try(name string, sc, cc *tls.Config) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	served := make(chan error, 1)
	go func() {
		raw, err := l.Accept()
		if err != nil {
			served <- err
			return
		}
		c := tls.Server(raw, sc)
		defer c.Close()
		if err := c.Handshake(); err != nil {
			served <- err
			return
		}
		b := make([]byte, 5)
		if _, err := io.ReadFull(c, b); err != nil {
			served <- err
			return
		}
		c.Write(b)
		served <- nil
	}()
	raw, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		panic(err)
	}
	c := tls.Client(raw, cc)
	c.SetDeadline(time.Now().Add(10 * time.Second))
	cerr := c.Handshake()
	echoed := false
	if cerr == nil {
		c.Write([]byte("hello"))
		b := make([]byte, 5)
		_, rerr := io.ReadFull(c, b)
		echoed = rerr == nil && string(b) == "hello"
	}
	c.Close()
	serr := <-served
	_, verify := errors.AsType[*tls.CertificateVerificationError](cerr)
	_, unknown := errors.AsType[x509.UnknownAuthorityError](cerr)
	fmt.Printf("%-24s echoed=%-5v server refused=%-5v client verify error=%-5v unknown authority=%v\n", name, echoed, serr != nil, verify, unknown)
}

func main() {
	now := time.Now()
	ca := issue("ca", nil, now.Add(-time.Hour), now.Add(time.Hour), 0)
	other := issue("other", nil, now.Add(-time.Hour), now.Add(time.Hour), 0)
	srv := issue("localhost", ca, now.Add(-time.Hour), now.Add(time.Hour), x509.ExtKeyUsageServerAuth)
	expired := issue("localhost", ca, now.Add(-2*time.Hour), now.Add(-time.Hour), x509.ExtKeyUsageServerAuth)
	client := issue("client", ca, now.Add(-time.Hour), now.Add(time.Hour), x509.ExtKeyUsageClientAuth)
	stranger := issue("client", other, now.Add(-time.Hour), now.Add(time.Hour), x509.ExtKeyUsageClientAuth)

	server := func(c *issued) *tls.Config { return &tls.Config{Certificates: []tls.Certificate{c.cert}} }
	trusting := func(p *x509.CertPool) *tls.Config { return &tls.Config{RootCAs: p, ServerName: "localhost"} }

	try("trusted", server(srv), trusting(pool(ca)))
	try("unknown authority", server(srv), trusting(pool(other)))
	try("empty roots", server(srv), trusting(x509.NewCertPool()))
	try("other name", server(srv), &tls.Config{RootCAs: pool(ca), ServerName: "elsewhere"})
	try("expired", server(expired), trusting(pool(ca)))
	try("skip verify", server(srv), &tls.Config{InsecureSkipVerify: true})

	pinned := trusting(pool(ca))
	pinned.VerifyPeerCertificate = func(raw [][]byte, chains [][]*x509.Certificate) error {
		return errors.New("pin mismatch")
	}
	try("peer check refuses", server(srv), pinned)
	checked := trusting(pool(ca))
	checked.VerifyConnection = func(s tls.ConnectionState) error {
		if s.PeerCertificates[0].Subject.CommonName != "someone else" {
			return errors.New("not who we want")
		}
		return nil
	}
	try("connection check", server(srv), checked)

	// A server's chain with a CA appended that signed none of it: the
	// verified chains end at the CA that did, and a pin on the other
	// does not hold.
	pinnedCA := issue("pinned", nil, now.Add(-time.Hour), now.Add(time.Hour), 0)
	padded := server(srv)
	padded.Certificates[0].Certificate = append(padded.Certificates[0].Certificate, pinnedCA.leaf.Raw)
	pin := trusting(pool(ca, pinnedCA))
	pin.VerifyPeerCertificate = func(raw [][]byte, chains [][]*x509.Certificate) error {
		for _, chain := range chains {
			for _, c := range chain {
				if c.Subject.CommonName == "pinned" {
					return nil
				}
			}
		}
		return errors.New("pinned CA not in the verified chains")
	}
	try("pin on a CA not signing", padded, pin)
	chains := trusting(pool(ca, pinnedCA))
	chains.VerifyConnection = func(s tls.ConnectionState) error {
		for _, chain := range s.VerifiedChains {
			fmt.Println("  verified chain ends at", chain[len(chain)-1].Subject.CommonName, "length", len(chain))
		}
		return nil
	}
	try("chains reported", padded, chains)

	mtls := server(srv)
	mtls.ClientAuth, mtls.ClientCAs = tls.RequireAndVerifyClientCert, pool(ca)
	withCert := func(c *issued) *tls.Config {
		cfg := trusting(pool(ca))
		cfg.Certificates = []tls.Certificate{c.cert}
		return cfg
	}
	try("mtls client trusted", mtls, withCert(client))
	try("mtls client stranger", mtls, withCert(stranger))
	try("mtls no client cert", mtls, trusting(pool(ca)))
	try("mtls server cert", mtls, withCert(srv))

	tls13 := server(srv)
	tls13.MinVersion = tls.VersionTLS13
	old := trusting(pool(ca))
	old.MaxVersion = tls.VersionTLS12
	try("versions apart", tls13, old)

	err := tls.Client(&net.TCPConn{}, &tls.Config{}).Handshake()
	fmt.Println("no server name:", err)
}
