//go:build gots_overlay

// Overlay for crypto/rand's fatal error: a panic.

package rand

func fatal(s string) { panic("fatal error: " + s) }
