//go:build gots_overlay

// Overlay for how deep the decoder skips nested types it does not decode
// into: Go's 10000 is sized for a goroutine's stack, which grows. A
// JavaScript stack does not, and a stream from outside nested a few
// thousand deep overflowed it and ended the program, the stack exhaustion
// of CVE-2024-34156; past 256 levels it is refused with an error.

package gob

var maxIgnoreNestingDepth = 256
