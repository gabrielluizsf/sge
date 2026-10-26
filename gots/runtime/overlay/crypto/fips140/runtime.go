//go:build gots_overlay

// Overlay for the FIPS 140 enforcement bypass Go keeps per goroutine:
// one for the program, as there is only one thread.

package fips140

var bypassed bool

func setBypass() { bypassed = true }

func isBypassed() bool { return bypassed }

func unsetBypass() { bypassed = false }
