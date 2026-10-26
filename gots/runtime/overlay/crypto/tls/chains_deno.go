//go:build gots_overlay

// The verified chains of a client's connection. node:tls tells whether
// the server's chain verifies, but not how: the chain it reports is the
// one the server sent, which may carry certificates that signed nothing
// (a CA a program pins, say, put there by a server certified by another).
// So Go builds the chains from the roots up, signatures checked, as its
// own handshake does: slow under gots, so only once something asks for
// them, the program's checks during the handshake or ConnectionState.

package tls

import (
	"crypto/x509"
	"errors"
	"syscall/js"
)

func (c *Conn) denoChains() error {
	p := c.pump()
	if !p.chainsPending {
		return nil
	}
	p.chainsPending = false
	roots, err := c.denoRoots()
	if err != nil {
		return err
	}
	certs := c.peerCertificates
	opts := x509.VerifyOptions{
		Roots:         roots,
		CurrentTime:   c.config.time(),
		DNSName:       c.config.ServerName,
		Intermediates: x509.NewCertPool(),
	}
	for _, cert := range certs[1:] {
		opts.Intermediates.AddCert(cert)
	}
	chains, err := certs[0].Verify(opts)
	if err != nil {
		return err
	}
	c.verifiedChains = chains
	return nil
}

// denoRoots are the roots the server's chain was verified against: RootCAs,
// or the runtime's, with what a system pool added to them.
func (c *Conn) denoRoots() (*x509.CertPool, error) {
	pool := c.config.RootCAs
	var extra []string
	if pool != nil {
		pems, system, err := x509.GotsPEMs(pool)
		if err != nil || !system {
			return pool, err
		}
		extra = pems
	}
	roots := x509.NewCertPool()
	for _, pem := range append(runtimeRoots(), extra...) {
		roots.AppendCertsFromPEM([]byte(pem))
	}
	if len(extra) == 0 && len(runtimeRoots()) == 0 {
		return nil, errors.New("tls: the runtime's roots cannot be read (Deno needs --allow-sys=ca)")
	}
	return roots, nil
}

var runtimeRootsPEM []string

// runtimeRoots are the roots node:tls verifies with by default, in PEM.
func runtimeRoots() []string {
	if len(runtimeRootsPEM) == 0 {
		// Read again until there are some: node:tls may not be loaded yet.
		v := js.Global().Get("gotsTls").Call("rootCertificates")
		runtimeRootsPEM = make([]string, v.Length())
		for i := range runtimeRootsPEM {
			runtimeRootsPEM[i] = v.Index(i).String()
		}
	}
	return runtimeRootsPEM
}
