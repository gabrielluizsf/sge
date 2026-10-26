//go:build gots_overlay

// A TCP listener of Deno's, as a netFD accepts on it.

package net

import (
	"os"
	"syscall/js"
	"time"
)

type denoListener struct {
	v        js.Value
	closed   bool
	deadline time.Time
}

// accept is the next connection, waiting for one until the deadline.
func (l *denoListener) accept() (*denoSock, error) {
	for {
		if l.closed {
			return nil, ErrClosed
		}
		r, err := awaitJS(l.v.Call("accept", waitMillis(l.deadline)))
		switch {
		case l.closed:
			// Deno keeps a listener open until a pending accept is done,
			// which may take one more connection: it is not ours.
			if err == nil && r.Type() == js.TypeObject {
				r.Call("close")
			}
			return nil, ErrClosed
		case err != nil:
			return nil, err
		case r.Type() == js.TypeString:
			if passed(l.deadline) {
				return nil, os.ErrDeadlineExceeded
			}
			continue
		}
		return &denoSock{v: r}, nil
	}
}

func (l *denoListener) setDeadline(t time.Time) {
	l.deadline = t
	l.v.Call("wake")
}

func (l *denoListener) close() error {
	if l.closed {
		return ErrClosed
	}
	l.closed = true
	l.v.Call("close")
	return nil
}
