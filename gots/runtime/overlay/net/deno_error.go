//go:build gots_overlay

// The errors Deno's sockets fail with, as the errnos Go programs test for.

package net

import (
	"errors"
	"syscall"
	"syscall/js"
)

// denoErrnos are the errnos of the Deno errors that have one.
var denoErrnos = map[string]syscall.Errno{
	"ConnectionRefused": syscall.ECONNREFUSED,
	"ConnectionReset":   syscall.ECONNRESET,
	"ConnectionAborted": syscall.ECONNABORTED,
	"NotConnected":      syscall.ENOTCONN,
	"AddrInUse":         syscall.EADDRINUSE,
	"AddrNotAvailable":  syscall.EADDRNOTAVAIL,
	"BrokenPipe":        syscall.EPIPE,
	"TimedOut":          syscall.ETIMEDOUT,
}

// denoError is the Go error for the error v a Deno call failed with: an
// errno where there is one; Deno's own message otherwise, which for a
// missing --allow-net says what to grant.
func denoError(v js.Value) error {
	name := v.Get("name").String()
	if e, ok := denoErrnos[name]; ok {
		return e
	}
	switch name {
	case "BadResource", "Interrupted":
		return ErrClosed
	}
	return errors.New(v.Get("message").String())
}
