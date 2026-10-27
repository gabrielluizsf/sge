//go:build gots_overlay

// A datagram socket of Deno's, UDP or Unix, as a netFD sends and
// receives on it.

package net

import (
	"os"
	"syscall/js"
	"time"
)

type denoDgram struct {
	v         js.Value
	closed    bool
	rdeadline time.Time
	wdeadline time.Time
	// peer is the address a connected socket, one a Dial made, sends to
	// and only takes datagrams from.
	peer *dgramAddr
	// unlink is the file of a socket bound to a path of gots's own choosing,
	// removed when it closes.
	unlink string
}

// dgramAddr is an address as the bridge has it: a host and a port, or a
// Unix socket's path and -1.
type dgramAddr struct {
	host string
	port int
}

func udpDgramAddr(a *UDPAddr) *dgramAddr { return &dgramAddr{a.IP.String(), a.Port} }

func (a dgramAddr) udp() *UDPAddr { return &UDPAddr{IP: ParseIP(a.host), Port: a.port} }

// same reports whether a and b are one address.
func (a dgramAddr) same(b *dgramAddr) bool {
	if a.port < 0 {
		return a.host == b.host && b.port < 0
	}
	return a.port == b.port && ParseIP(a.host).Equal(ParseIP(b.host))
}

// receive reads the next datagram into b, cut to fit as Go's are, and
// says who sent it, waiting until the read deadline.
func (d *denoDgram) receive(b []byte) (int, dgramAddr, error) {
	for {
		if d.closed {
			return 0, dgramAddr{}, ErrClosed
		}
		r, err := awaitJS(d.v.Call("receive", waitMillis(d.rdeadline)))
		switch {
		case d.closed:
			return 0, dgramAddr{}, ErrClosed
		case err != nil:
			return 0, dgramAddr{}, err
		case r.Type() == js.TypeString:
			if passed(d.rdeadline) {
				return 0, dgramAddr{}, os.ErrDeadlineExceeded
			}
			continue
		}
		from := dgramAddr{r.Index(1).String(), r.Index(2).Int()}
		if d.peer != nil && !from.same(d.peer) {
			continue // a connected socket hears only its peer
		}
		return js.CopyBytesToGo(b, r.Index(0)), from, nil
	}
}

// send sends b to to as one datagram.
func (d *denoDgram) send(b []byte, to *dgramAddr) (int, error) {
	if d.closed {
		return 0, ErrClosed
	}
	if passed(d.wdeadline) {
		return 0, os.ErrDeadlineExceeded
	}
	a := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(a, b)
	n, err := awaitJS(d.v.Call("send", a, to.host, to.port))
	if err != nil {
		return 0, err
	}
	return n.Int(), nil
}

func (d *denoDgram) setReadDeadline(t time.Time) {
	d.rdeadline = t
	d.v.Call("wake")
}

func (d *denoDgram) close() error {
	if d.closed {
		return ErrClosed
	}
	d.closed = true
	d.v.Call("close")
	if d.unlink != "" {
		os.Remove(d.unlink)
	}
	return nil
}
