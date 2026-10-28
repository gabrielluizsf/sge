//go:build gots_overlay

// Overlay for data-independent timing, a CPU mode JavaScript has no way
// to ask for: it is never on.

package subtle

func setDITEnabled() bool { return false }

func setDITDisabled() {}
