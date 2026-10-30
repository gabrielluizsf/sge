// tlsfull: crypto/tls as Go has it, over loopback: keying material,
// stapled OCSP, the versions and suites a config allows, client
// certificates required and checked, session resumption, and certificates
// picked with no name asked for.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"time"
)

var now = time.Now()

func main() {
	ca, caKey := cert("ca", nil, nil, true)
	srv, srvKey := cert("server", ca, caKey, false)
	cli, cliKey := cert("client", ca, caKey, false)
	roots := x509.NewCertPool()
	roots.AddCert(ca)
	serverCert := tls.Certificate{Certificate: [][]byte{srv.Raw}, PrivateKey: srvKey, OCSPStaple: []byte("stapled")}
	clientCert := tls.Certificate{Certificate: [][]byte{cli.Raw}, PrivateKey: cliKey}

	cache := tls.NewLRUClientSessionCache(4)
	base := func() *tls.Config {
		return &tls.Config{ServerName: "server", RootCAs: roots, ClientSessionCache: cache, Time: func() time.Time { return now }}
	}
	server := &tls.Config{Certificates: []tls.Certificate{serverCert}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: roots, Time: func() time.Time { return now }}

	c := base()
	c.Certificates = []tls.Certificate{clientCert}
	report("full", c, server)
	report("resumed", c, server)

	// A client with no certificate is refused.
	report("no client cert", base(), server)

	// TLS 1.3 only, against a client of TLS 1.2 only.
	s13 := server.Clone()
	s13.MinVersion = tls.VersionTLS13
	s13.ClientAuth = tls.NoClientCert
	c12 := base()
	c12.MaxVersion = tls.VersionTLS12
	report("1.3 only", c12, s13)

	// TLS 1.2 with the suites the config names.
	s12 := server.Clone()
	s12.ClientAuth = tls.NoClientCert
	s12.MaxVersion = tls.VersionTLS12
	c12.CipherSuites = []uint16{tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305}
	report("1.2 chacha", c12, s12)
	c12.CipherSuites = []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}
	report("1.2 aes", c12, s12)

	// No name asked for: GetCertificate picks.
	sn := &tls.Config{GetCertificate: func(h *tls.ClientHelloInfo) (*tls.Certificate, error) {
		fmt.Println("asked for", fmt.Sprintf("%q", h.ServerName))
		return &serverCert, nil
	}}
	insecure := &tls.Config{InsecureSkipVerify: true}
	report("no name", insecure, sn)
}

// report runs a handshake and an echo, and prints what was agreed.
func report(name string, client, server *tls.Config) {
	l, err := tls.Listen("tcp", "127.0.0.1:0", server)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	done := make(chan string, 1)
	go func() {
		c, err := l.Accept()
		if err != nil {
			done <- err.Error()
			return
		}
		defer c.Close()
		tc := c.(*tls.Conn)
		if err := tc.Handshake(); err != nil {
			done <- "server: " + err.Error()
			return
		}
		st := tc.ConnectionState()
		ekm, _ := st.ExportKeyingMaterial("label", []byte("ctx"), 16)
		b := make([]byte, 5)
		io.ReadFull(tc, b)
		tc.Write(b)
		done <- fmt.Sprint("server peers ", len(st.PeerCertificates), " chains ", len(st.VerifiedChains), " ekm ", hex.EncodeToString(ekm) != "")
	}()
	c, err := tls.Dial("tcp", l.Addr().String(), client)
	if err != nil {
		fmt.Println(name, "dial:", err)
		fmt.Println(name, <-done)
		return
	}
	defer c.Close()
	if _, err := c.Write([]byte("hello")); err != nil {
		fmt.Println(name, "write:", err)
	}
	b := make([]byte, 5)
	_, err = io.ReadFull(c, b)
	st := c.ConnectionState()
	ekm, kerr := st.ExportKeyingMaterial("label", []byte("ctx"), 16)
	fmt.Println(name, string(b), err, tls.VersionName(st.Version), suite(st),
		"resumed", st.DidResume, "ocsp", string(st.OCSPResponse), "chains", len(st.VerifiedChains), "ekm", len(ekm), kerr)
	fmt.Println(name, <-done)
}

// suite is the connection's cipher suite. TLS 1.3 picks AES-GCM or
// ChaCha20 as the machine's CPU has AES hardware or not, which Go asks of
// the CPU and a transpiled program cannot, so either is printed the same;
// TLS 1.2 uses the one its config names.
func suite(st tls.ConnectionState) string {
	switch st.CipherSuite {
	case tls.TLS_AES_128_GCM_SHA256, tls.TLS_CHACHA20_POLY1305_SHA256:
		return "AES-128-GCM or CHACHA20-POLY1305"
	}
	return tls.CipherSuiteName(st.CipherSuite)
}

func cert(name string, parent *x509.Certificate, parentKey *ecdsa.PrivateKey, isCA bool) (*x509.Certificate, *ecdsa.PrivateKey) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	t := &x509.Certificate{
		SerialNumber: big.NewInt(int64(len(name))), Subject: pkix.Name{CommonName: name},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), DNSNames: []string{name},
		IsCA: isCA, BasicConstraintsValid: true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	if parent == nil {
		parent, parentKey = t, key
	}
	der, err := x509.CreateCertificate(rand.Reader, t, parent, &key.PublicKey, parentKey)
	if err != nil {
		panic(err)
	}
	c, _ := x509.ParseCertificate(der)
	return c, key
}
