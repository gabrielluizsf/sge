//go:build gots_overlay

// The connection a Conn runs over under gots: Go's own, whose bytes two
// goroutines carry to and from node:tls, which speaks the protocol.

package tls

import (
	"io"
	"net"
	"os"
	"syscall/js"
	"time"
)

type pumpConn struct {
	net.Conn // the connection TLS runs over

	side                 js.Value      // the bridge's Side, once the handshake starts
	handshaking          chan struct{} // closed once the handshake is over
	ending               bool          // close_notify is on its way
	closeAfterOut        bool          // Close waits for it to be sent
	chainsPending        bool          // see denoChains
	closed               bool
	rdeadline, wdeadline time.Time
}

// Write drops what Go's own record layer writes: the alerts it sends on
// its way out of a failed check, in the clear, as it never has keys.
// node:tls sends what is sent.
func (p *pumpConn) Write(b []byte) (int, error) { return len(b), nil }

// Close closes the connection, once close_notify, if on its way, is sent
// (at most 5 seconds, as Go waits): Close itself does not wait, as code
// holding a lock closes connections, and a lock of gots's cannot be
// waited for.
func (p *pumpConn) Close() error {
	if p.closed {
		return net.ErrClosed
	}
	p.closed = true
	if p.started() && p.ending {
		p.closeAfterOut = true
		time.AfterFunc(5*time.Second, func() { p.Conn.Close() })
		return nil
	}
	if p.started() {
		p.side.Call("destroy")
	}
	return p.Conn.Close()
}

// end has close_notify sent, without waiting for it to be.
func (p *pumpConn) end() {
	if p.started() && !p.ending && !p.closed {
		p.ending = true
		p.side.Call("end")
	}
}

func (p *pumpConn) started() bool { return p.side.Type() == js.TypeObject }

// start opens the Side, a server or a client with opts, and its pumps.
func (p *pumpConn) start(server bool, opts map[string]any) error {
	p.side = js.Global().Get("gotsTls").Call("side")
	if _, err := awaitJS(p.side.Call("open", server, opts)); err != nil {
		return err
	}
	go p.pumpIn()
	go p.pumpOut()
	return nil
}

// pumpIn feeds TLS what the peer sends, until it sends no more.
func (p *pumpConn) pumpIn() {
	buf := make([]byte, 32<<10)
	for {
		n, err := p.Conn.Read(buf)
		if n > 0 {
			a := js.Global().Get("Uint8Array").New(n)
			js.CopyBytesToJS(a, buf[:n])
			p.side.Call("feed", a)
		}
		if err != nil {
			if !p.closed {
				p.side.Call("feedEnd")
			}
			return
		}
	}
}

// pumpOut sends the peer what TLS sends, until it sends no more.
func (p *pumpConn) pumpOut() {
	for {
		v, err := awaitJS(p.side.Call("takeOut"))
		if err != nil || v.IsNull() {
			if p.closeAfterOut {
				p.Conn.Close()
			}
			return
		}
		b := make([]byte, v.Length())
		js.CopyBytesToGo(b, v)
		if _, err := p.Conn.Write(b); err != nil {
			p.side.Call("destroy")
			return
		}
	}
}

// read reads into b what the peer sent, waiting for some until the read
// deadline. A read that runs out of time gives up nothing.
func (p *pumpConn) read(b []byte) (int, error) {
	for {
		if p.closed {
			return 0, net.ErrClosed
		}
		r, err := awaitJS(p.side.Call("read", len(b), waitMillis(p.rdeadline)))
		switch {
		case p.closed:
			return 0, net.ErrClosed
		case err != nil:
			return 0, err
		case r.IsNull():
			return 0, io.EOF
		case r.Type() == js.TypeString:
			// Out of time, or woken by a new deadline.
			if passed(p.rdeadline) {
				return 0, os.ErrDeadlineExceeded
			}
			continue
		}
		return js.CopyBytesToGo(b, r), nil
	}
}

// write writes all of b until the write deadline, as net's sockets do.
func (p *pumpConn) write(b []byte) (int, error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	if passed(p.wdeadline) {
		return 0, os.ErrDeadlineExceeded
	}
	a := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(a, b)
	w := p.side.Call("startWrite", a)
	for {
		if _, err := awaitJS(w.Call("wait", waitMillis(p.wdeadline))); err != nil {
			return 0, err
		}
		switch {
		case p.closed:
			return w.Get("n").Int(), net.ErrClosed
		case w.Get("done").Bool():
			if e := w.Get("err"); !e.IsNull() {
				return w.Get("n").Int(), jsError(e)
			}
			return w.Get("n").Int(), nil
		case passed(p.wdeadline):
			return w.Get("n").Int(), os.ErrDeadlineExceeded
		}
	}
}

// waitMillis is how long a call may wait for deadline d: -1 for no
// deadline, 0 for one passed, otherwise at least a millisecond.
func waitMillis(d time.Time) int {
	if d.IsZero() {
		return -1
	}
	left := time.Until(d)
	if left <= 0 {
		return 0
	}
	return int((left + time.Millisecond - 1) / time.Millisecond)
}

func passed(d time.Time) bool { return !d.IsZero() && !time.Now().Before(d) }
