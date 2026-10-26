// tlsecho: a TLS server and client over loopback, with a certificate made
// here, the client trusting only it.
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	"net"
	"time"
)

func main() {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "gots test"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour),
		DNSNames: []string{"localhost"}, IsCA: true, BasicConstraintsValid: true,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	leaf, _ := x509.ParseCertificate(der)
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	cert := tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}

	l, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}, NextProtos: []string{"echo/1"}})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	go func() {
		c, err := l.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		io.Copy(c, c)
	}()

	c, err := tls.Dial("tcp", l.Addr().String(), &tls.Config{RootCAs: pool, ServerName: "localhost", NextProtos: []string{"echo/1"}})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	fmt.Fprintf(c, "hello over tls")
	b := make([]byte, 64)
	n, err := c.Read(b)
	fmt.Println(string(b[:n]), err)
	st := c.ConnectionState()
	fmt.Println("version", tls.VersionName(st.Version), "alpn", st.NegotiatedProtocol, "server", st.ServerName)
	fmt.Println("peer", st.PeerCertificates[0].Subject.CommonName, "chains", len(st.VerifiedChains))
	_ = net.IPv4len
}
