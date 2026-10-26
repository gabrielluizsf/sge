//go:build gots_overlay

// What a handshake of node:tls's is asked to agree on, and what Go checks
// it did agree on: node:tls may settle on a version it was told not to
// (a server of Deno's has taken TLS 1.2 though told 1.3 at least), so
// nothing it reports is taken on trust where Go can look.

package tls

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall/js"
)

var nodeVersions = map[uint16]string{VersionTLS12: "TLSv1.2", VersionTLS13: "TLSv1.3"}

// denoVersions are the versions c allows that node:tls speaks, 1.2 and
// 1.3: fewer than Go's at most, never more.
func denoVersions(c *Config) (lo, hi uint16, err error) {
	lo, hi = c.MinVersion, c.MaxVersion
	if lo < VersionTLS12 {
		lo = VersionTLS12
	}
	if hi == 0 || hi > VersionTLS13 {
		hi = VersionTLS13
	}
	if lo > hi {
		return 0, 0, errors.New("tls: no supported versions satisfy MinVersion and MaxVersion")
	}
	return lo, hi, nil
}

// versionOpts asks node:tls for the versions from lo to hi.
func versionOpts(opts map[string]any, lo, hi uint16) {
	opts["minVersion"] = nodeVersions[lo]
	opts["maxVersion"] = nodeVersions[hi]
}

// denoStep waits for the handshake's next step: its Info, or for a server
// the name it is to pick a certificate for. It waits until the read
// deadline or ctx's, as Go's handshake reads.
func (c *Conn) denoStep(ctx context.Context) (js.Value, error) {
	p := c.pump()
	for {
		deadline := p.rdeadline
		if d, ok := ctx.Deadline(); ok && (deadline.IsZero() || d.Before(deadline)) {
			deadline = d
		}
		r, err := awaitJS(p.side.Call("step", waitMillis(deadline)))
		if err != nil {
			return js.Undefined(), err
		}
		if r.Type() != js.TypeString {
			return r, nil
		}
		if err := ctx.Err(); err != nil {
			return js.Undefined(), err
		}
		if passed(deadline) {
			return js.Undefined(), os.ErrDeadlineExceeded
		}
	}
}

// denoAgreed takes in what info says the handshake agreed on, refusing a
// version outside lo to hi.
func (c *Conn) denoAgreed(info js.Value, lo, hi uint16) error {
	name := info.Get("version").String()
	var vers uint16
	for v, n := range nodeVersions {
		if n == name {
			vers = v
		}
	}
	if vers == 0 || vers < lo || vers > hi {
		return fmt.Errorf("tls: the handshake settled on %q, outside the versions configured", name)
	}
	c.vers, c.haveVers = vers, true
	cipher := info.Get("cipher").String()
	for _, s := range append(CipherSuites(), InsecureCipherSuites()...) {
		if s.Name == cipher {
			c.cipherSuite = s.ID
		}
	}
	c.ekm = noEKM
	return nil
}

// peerChain is the peer's certificates info holds, its own first.
func peerChain(info js.Value) [][]byte {
	chain := info.Get("chain")
	raw := make([][]byte, chain.Length())
	for i := range raw {
		v := chain.Index(i)
		raw[i] = make([]byte, v.Length())
		js.CopyBytesToGo(raw[i], v)
	}
	return raw
}

// fail ends a handshake that went wrong past node:tls: nothing more is
// sent, and err is what Handshake gives.
func (c *Conn) fail(err error) error {
	c.pump().side.Call("destroy")
	return err
}
