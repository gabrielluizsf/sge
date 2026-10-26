//go:build gots_overlay

// A UDP socket of Deno's, as a netFD sends and receives on it.

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
	// peer is the address a connected socket, one DialUDP made, sends to
	// and only takes datagrams from.
	peer *UDPAddr
}

// receive reads the next datagram into b, cut to fit as Go's are, and
// says who sent it, waiting until the read deadline.
func (d *denoDgram) receive(b []byte) (int, *UDPAddr, error) {
	for {
		if d.closed {
			return 0, nil, ErrClosed
		}
		r, err := awaitJS(d.v.Call("receive", waitMillis(d.rdeadline)))
		switch {
		case d.closed:
			return 0, nil, ErrClosed
		case err != nil:
			return 0, nil, err
		case r.Type() == js.TypeString:
			if passed(d.rdeadline) {
				return 0, nil, os.ErrDeadlineExceeded
			}
			continue
		}
		from := &UDPAddr{IP: ParseIP(r.Index(1).String()), Port: r.Index(2).Int()}
		if d.peer != nil && (!from.IP.Equal(d.peer.IP) || from.Port != d.peer.Port) {
			continue // a connected socket hears only its peer
		}
		return js.CopyBytesToGo(b, r.Index(0)), from, nil
	}
}

// send sends b to to as one datagram.
func (d *denoDgram) send(b []byte, to *UDPAddr) (int, error) {
	if d.closed {
		return 0, ErrClosed
	}
	if passed(d.wdeadline) {
		return 0, os.ErrDeadlineExceeded
	}
	a := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(a, b)
	n, err := awaitJS(d.v.Call("send", a, to.IP.String(), to.Port))
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
	return nil
}
