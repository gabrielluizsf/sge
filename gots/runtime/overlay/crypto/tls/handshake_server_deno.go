//go:build gots_overlay

// A server's handshake on Deno. node:tls does not check a client's
// certificate as a server, whatever it is told (it has taken one signed
// by nobody), so it only asks for one, and Go checks it as its own
// handshake does, with processCertsFromClient.

package tls

import (
	"context"
	"errors"
	"slices"
	"syscall/js"
)

func (c *Conn) denoServerHandshake(ctx context.Context) error {
	if c.config == nil {
		c.config = defaultConfig()
	}
	cfg := c.config
	lo, hi, err := denoVersions(cfg)
	if err != nil {
		return err
	}
	// A certificate is picked by the name asked for when there is a
	// choice; GetConfigForClient may then ask for a client's.
	pick := cfg.GetCertificate != nil || cfg.GetConfigForClient != nil || len(cfg.Certificates) > 1
	opts := map[string]any{
		"requestCert":        cfg.ClientAuth != NoClientCert || cfg.GetConfigForClient != nil,
		"rejectUnauthorized": false,
		"sni":                pick,
	}
	versionOpts(opts, lo, hi)
	if len(cfg.NextProtos) > 0 {
		opts["ALPNProtocols"] = anys(cfg.NextProtos)
	}
	if len(cfg.Certificates) > 0 {
		chain, key, err := pemPair(&cfg.Certificates[0])
		if err != nil {
			return err
		}
		opts["cert"], opts["key"] = chain, key
	} else if !pick {
		return errors.New("tls: no certificates configured")
	}

	p := c.pump()
	if err := p.start(true, opts); err != nil {
		return err
	}
	var info js.Value
	for {
		r, err := c.denoStep(ctx)
		if err != nil {
			return c.fail(err)
		}
		if sni := r.Get("sni"); sni.Type() == js.TypeString {
			if err := c.denoPick(ctx, sni.String()); err != nil {
				p.side.Call("refuseSni", "tls: no certificate for the name asked")
				return c.fail(err)
			}
			continue
		}
		info = r
		break
	}
	// GetConfigForClient may have put another config in c's.
	cfg = c.config
	if lo, hi, err = denoVersions(cfg); err != nil {
		return c.fail(err)
	}
	if err := c.denoAgreed(info, lo, hi); err != nil {
		return c.fail(err)
	}
	c.serverName = info.Get("serverName").String()
	alpn := info.Get("alpn").String()
	if alpn != "" && !slices.Contains(cfg.NextProtos, alpn) {
		return c.fail(errors.New("tls: the handshake settled on an application protocol not configured"))
	}
	c.clientProtocol = alpn
	if cfg.ClientAuth != NoClientCert {
		if err := c.processCertsFromClient(Certificate{Certificate: peerChain(info)}); err != nil {
			return c.fail(err)
		}
	}
	if f := cfg.VerifyConnection; f != nil {
		if err := f(c.connectionStateLocked()); err != nil {
			return c.fail(err)
		}
	}
	c.isHandshakeComplete.Store(true)
	return nil
}

// denoPick answers node:tls's question of which certificate to show a
// client asking for name, as Go's server picks it.
func (c *Conn) denoPick(ctx context.Context, name string) error {
	hello := c.denoHello(ctx, name)
	if f := c.config.GetConfigForClient; f != nil {
		nc, err := f(hello)
		if err != nil {
			return err
		}
		if nc != nil {
			c.config = nc
			hello.config = nc
		}
	}
	cert, err := c.config.getCertificate(hello)
	if err != nil {
		return err
	}
	chain, key, err := pemPair(cert)
	if err != nil {
		return err
	}
	_, err = awaitJS(c.pump().side.Call("answerSni", chain, key))
	return err
}

// denoHello is the ClientHelloInfo of a client asking for name. node:tls
// tells the name alone: the rest is what clients offer nowadays, so a
// certificate is picked by its name and key as Go picks it.
func (c *Conn) denoHello(ctx context.Context, name string) *ClientHelloInfo {
	return &ClientHelloInfo{
		ServerName:        name,
		SupportedVersions: []uint16{VersionTLS13, VersionTLS12},
		CipherSuites: []uint16{
			TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256,
			TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		SupportedCurves: []CurveID{X25519, CurveP256, CurveP384},
		SupportedPoints: []uint8{pointFormatUncompressed},
		SignatureSchemes: []SignatureScheme{
			ECDSAWithP256AndSHA256, ECDSAWithP384AndSHA384, Ed25519,
			PSSWithSHA256, PSSWithSHA384, PSSWithSHA512,
			PKCS1WithSHA256, PKCS1WithSHA384, PKCS1WithSHA512,
		},
		SupportedProtos: c.config.NextProtos,
		Conn:            c.pump().Conn,
		config:          c.config,
		ctx:             ctx,
	}
}
