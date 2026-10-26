//go:build gots_overlay

// Overlay for the random numbers net takes from the runtime, to shuffle
// DNS answers: they come from the platform's crypto.getRandomValues.

package net

import "syscall/js"

var entropy js.Value

func runtime_rand() uint64 {
	if entropy.IsUndefined() {
		entropy = js.Global().Get("Uint32Array").New(2)
	}
	js.Global().Get("crypto").Call("getRandomValues", entropy)
	return uint64(entropy.Index(0).Int())<<32 | uint64(entropy.Index(1).Int())
}
