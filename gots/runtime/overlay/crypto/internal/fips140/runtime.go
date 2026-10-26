//go:build gots_overlay

// Overlay for what crypto/internal/fips140 takes from the runtime: a
// fatal error is a panic, and the service indicator, which Go keeps per
// goroutine, is one for the program, as there is only one thread.

package fips140

var indicator uint8

func fatal(s string) { panic("fatal error: " + s) }

func getIndicator() uint8 { return indicator }

func setIndicator(v uint8) { indicator = v }
