//go:build gots_overlay

// Overlay for the markers BoringCrypto builds look for in a binary:
// there is none here.

package sig

func BoringCrypto() {}

func FIPSOnly() {}

func StandardCrypto() {}
