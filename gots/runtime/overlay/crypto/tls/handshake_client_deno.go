//go:build gots_overlay

// A client's handshake on Deno. node:tls checks the server's certificate
// against RootCAs, which take the place of the runtime's roots as they
// take the place of the system's in Go, or against the runtime's roots
// when RootCAs is nil. It does not stop the handshake for it, which would
// tell Go no more than that the chain failed: Go stops it, before a byte
// of the program's is sent, with the error its own check would give, and
// checks the name and the dates again itself. Only a client showing a
// certificate has node:tls stop, so it is never shown to a server that
// failed.

package tls

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"
	"syscall/js"
	"time"
)

func (c *Conn) denoClientHandshake(ctx context.Context) error {
	if c.config == nil {
		c.config = defaultConfig()
	}
	cfg := c.config
	if cfg.ServerName == "" && !cfg.InsecureSkipVerify {
		return errors.New("tls: either ServerName or InsecureSkipVerify must be specified in the tls.Config")
	}
	lo, hi, err := denoVersions(cfg)
	if err != nil {
		return err
	}
	cert, err := c.denoClientCert(ctx)
	if err != nil {
		return err
	}
	shows := cert != nil && len(cert.Certificate) > 0
	opts := map[string]any{"rejectUnauthorized": !cfg.InsecureSkipVerify && shows}
	versionOpts(opts, lo, hi)
	name := strings.TrimSuffix(cfg.ServerName, ".")
	switch {
	case name == "":
		opts["servername"] = "" // no SNI, as Go sends none
	case net.ParseIP(strings.Trim(name, "[]")) != nil:
		opts["host"] = strings.Trim(name, "[]") // checked as an address, sent as no SNI
	default:
		opts["host"] = name
		opts["servername"] = name
	}
	if !cfg.InsecureSkipVerify && cfg.RootCAs != nil {
		pems, system, err := x509.GotsPEMs(cfg.RootCAs)
		if err != nil {
			return err
		}
		opts["ca"] = anys(pems)
		opts["systemRoots"] = system
	}
	if len(cfg.NextProtos) > 0 {
		opts["ALPNProtocols"] = anys(cfg.NextProtos)
	}
	if shows {
		chain, key, err := pemPair(cert)
		if err != nil {
			return err
		}
		opts["cert"], opts["key"] = chain, key
	}

	p := c.pump()
	if err := p.start(false, opts); err != nil {
		return err
	}
	info, err := c.denoStep(ctx)
	if err != nil {
		if !cfg.InsecureSkipVerify {
			err = verifyError(err)
		}
		return c.fail(err)
	}
	if err := c.denoAgreed(info, lo, hi); err != nil {
		return c.fail(err)
	}
	alpn := info.Get("alpn").String()
	if alpn != "" && !slices.Contains(cfg.NextProtos, alpn) {
		return c.fail(errors.New("tls: server selected unadvertised ALPN protocol"))
	}
	c.clientProtocol = alpn
	c.serverName = cfg.ServerName
	if err := c.denoServerCerts(info); err != nil {
		return c.fail(err)
	}
	c.isHandshakeComplete.Store(true)
	return nil
}

// denoClientCert is the certificate the client shows, if asked. node:tls
// takes it before the handshake starts, so GetClientCertificate is asked
// before what the server accepts is known.
func (c *Conn) denoClientCert(ctx context.Context) (*Certificate, error) {
	if f := c.config.GetClientCertificate; f != nil {
		return f(&CertificateRequestInfo{Version: VersionTLS13, ctx: ctx})
	}
	if len(c.config.Certificates) > 0 {
		return &c.config.Certificates[0], nil
	}
	return nil, nil
}

// denoServerCerts takes in the server's chain as verifyServerCertificate
// does, but for the verifying, which node:tls did. The verified chains
// are Go's to work out (see denoChains), when asked for.
func (c *Conn) denoServerCerts(info js.Value) error {
	raw := peerChain(info)
	if len(raw) == 0 {
		return errors.New("tls: server sent no certificate")
	}
	certs := make([]*x509.Certificate, len(raw))
	for i, der := range raw {
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return errors.New("tls: failed to parse certificate from server: " + err.Error())
		}
		if cert.PublicKeyAlgorithm == x509.RSA {
			if max, ok := checkKeySize(cert.PublicKey.(*rsa.PublicKey).N.BitLen()); !ok {
				return fmt.Errorf("tls: server sent certificate containing RSA key larger than %d bits", max)
			}
		}
		certs[i] = cert
	}
	switch certs[0].PublicKey.(type) {
	case *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey:
	default:
		return fmt.Errorf("tls: server's certificate contains an unsupported type of public key: %T", certs[0].PublicKey)
	}
	if !c.config.InsecureSkipVerify {
		if err := c.denoTrusted(certs, info); err != nil {
			return &CertificateVerificationError{UnverifiedCertificates: certs, Err: err}
		}
	}
	c.peerCertificates = certs
	if !c.config.InsecureSkipVerify {
		c.pump().chainsPending = true
		if c.config.VerifyPeerCertificate != nil || c.config.VerifyConnection != nil {
			if err := c.denoChains(); err != nil {
				return &CertificateVerificationError{UnverifiedCertificates: certs, Err: err}
			}
		}
	}
	if f := c.config.VerifyPeerCertificate; f != nil {
		if err := f(raw, c.verifiedChains); err != nil {
			return err
		}
	}
	if f := c.config.VerifyConnection; f != nil {
		if err := f(c.connectionStateLocked()); err != nil {
			return err
		}
	}
	return nil
}

// denoTrusted is why the server's chain is not to be trusted, or nil: as
// node:tls found, explained as Go's check would, and the name and dates
// as Go checks them.
func (c *Conn) denoTrusted(certs []*x509.Certificate, info js.Value) error {
	now := c.config.time()
	for _, cert := range certs {
		if now.Before(cert.NotBefore) {
			return x509.CertificateInvalidError{Cert: cert, Reason: x509.Expired,
				Detail: fmt.Sprintf("current time %s is before %s", now.Format(time.RFC3339), cert.NotBefore.Format(time.RFC3339))}
		}
		if now.After(cert.NotAfter) {
			return x509.CertificateInvalidError{Cert: cert, Reason: x509.Expired,
				Detail: fmt.Sprintf("current time %s is after %s", now.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339))}
		}
	}
	if !info.Get("authorized").Bool() && info.Get("authError").String() != "ERR_TLS_CERT_ALTNAME_INVALID" {
		return x509.UnknownAuthorityError{Cert: certs[0]}
	}
	if err := certs[0].VerifyHostname(c.config.ServerName); err != nil {
		return err
	}
	if !info.Get("authorized").Bool() {
		// Refused for its name, though Go finds it fine.
		return fmt.Errorf("x509: certificate is not valid for %q", c.config.ServerName)
	}
	return nil
}
