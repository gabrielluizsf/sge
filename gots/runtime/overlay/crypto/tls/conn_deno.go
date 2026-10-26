//go:build gots_overlay

// Overlay for Conn on Deno: node:tls speaks the protocol (see pumpConn),
// and a Conn reads and writes what it decrypts and encrypts. What the
// handshake agreed on fills the fields Go's own would, so ConnectionState,
// VerifyHostname and Close work on as they are.

package tls

import (
	"errors"
	"net"
	"time"
)

func Server(conn net.Conn, config *Config) *Conn {
	c := &Conn{conn: &pumpConn{Conn: conn}, config: config}
	c.handshakeFn = c.denoServerHandshake
	return c
}

func Client(conn net.Conn, config *Config) *Conn {
	c := &Conn{conn: &pumpConn{Conn: conn}, config: config, isClient: true}
	c.handshakeFn = c.denoClientHandshake
	return c
}

func (c *Conn) pump() *pumpConn { return c.conn.(*pumpConn) }

func (c *Conn) NetConn() net.Conn { return c.pump().Conn }

// The deadlines are the reads' and writes' of what TLS carries: the pumps
// keep reading and writing the connection under it meanwhile.

func (c *Conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	return c.SetWriteDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	p := c.pump()
	p.rdeadline = t
	if p.started() {
		p.side.Call("wake")
	}
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	p := c.pump()
	p.wdeadline = t
	if p.started() {
		p.side.Call("wakeWriters")
	}
	return nil
}

func (c *Conn) Read(b []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, nil
	}
	return c.pump().read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	// As Go's: a Write may not start once Close has.
	for {
		x := c.activeCall.Load()
		if x&1 != 0 {
			return 0, net.ErrClosed
		}
		if c.activeCall.CompareAndSwap(x, x+2) {
			break
		}
	}
	defer c.activeCall.Add(-2)
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	if c.closeNotifySent {
		return 0, errShutdown
	}
	return c.pump().write(b)
}

// closeNotify has close_notify sent, without waiting: see pumpConn.Close.
func (c *Conn) closeNotify() error {
	if !c.closeNotifySent {
		c.closeNotifySent = true
		c.pump().end()
	}
	return nil
}

var errNoEKM = errors.New("crypto/tls: ExportKeyingMaterial is not available under gots: node:tls has no keying material exporter")

func noEKM(string, []byte, int) ([]byte, error) { return nil, errNoEKM }
