//go:build gots_overlay

// Overlay for the system's randomness on js: the platform's
// crypto.getRandomValues, a cryptographically secure generator, as Go's
// own js port uses; and a fatal error is a panic.

package sysrand

import "syscall/js"

func fatal(s string) { panic("fatal error: " + s) }

func getRandomValues(r []byte) {
	a := js.Global().Get("Uint8Array").New(len(r))
	js.Global().Get("crypto").Call("getRandomValues", a)
	js.CopyBytesToGo(r, a)
}
