//go:build gots_overlay

// A TCP connection of Deno's, as a netFD reads and writes it.

package net

import (
	"io"
	"os"
	"syscall/js"
	"time"
)

type denoSock struct {
	v                    js.Value
	closed               bool
	rdeadline, wdeadline time.Time
}

// read reads into b what has arrived, waiting for some until the read
// deadline. A read that runs out of time gives up nothing: what arrives
// later is the next read's.
func (s *denoSock) read(b []byte) (int, error) {
	if s.closed {
		return 0, ErrClosed
	}
	if len(b) == 0 {
		return 0, nil
	}
	for {
		r, err := awaitJS(s.v.Call("read", len(b), waitMillis(s.rdeadline)))
		switch {
		case s.closed:
			return 0, ErrClosed
		case err != nil:
			return 0, err
		case r.IsNull():
			return 0, io.EOF
		case r.Type() == js.TypeString:
			// Out of time, or woken by a new deadline.
			if passed(s.rdeadline) {
				return 0, os.ErrDeadlineExceeded
			}
			continue
		}
		return js.CopyBytesToGo(b, r), nil
	}
}

// write writes all of b. The write deadline is held only before writing:
// once handed to Deno a write goes through.
func (s *denoSock) write(b []byte) (int, error) {
	if s.closed {
		return 0, ErrClosed
	}
	if passed(s.wdeadline) {
		return 0, os.ErrDeadlineExceeded
	}
	a := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(a, b)
	n, err := awaitJS(s.v.Call("write", a))
	if err != nil {
		if s.closed {
			return 0, ErrClosed
		}
		return 0, err
	}
	return n.Int(), nil
}

func (s *denoSock) setReadDeadline(t time.Time) {
	s.rdeadline = t
	s.v.Call("wake")
}

func (s *denoSock) setWriteDeadline(t time.Time) { s.wdeadline = t }

func (s *denoSock) closeWrite() error {
	if s.closed {
		return ErrClosed
	}
	if e := s.v.Call("closeWrite"); !e.IsNull() {
		return denoError(e)
	}
	return nil
}

func (s *denoSock) close() error {
	if s.closed {
		return ErrClosed
	}
	s.closed = true
	s.v.Call("close")
	return nil
}

// tcpAddr is the address the bridge gives as [host, port].
func tcpAddr(a js.Value) *TCPAddr {
	return &TCPAddr{IP: ParseIP(a.Index(0).String()), Port: a.Index(1).Int()}
}
