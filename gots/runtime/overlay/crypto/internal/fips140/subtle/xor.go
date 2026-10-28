//go:build gots_overlay

// Overlay for XORBytes, which Go hands to code that reads the bytes a
// machine word at a time through unsafe, aligned or not by their address:
// here they go a byte at a time.

package subtle

import "crypto/internal/fips140/alias"

func XORBytes(dst, x, y []byte) int {
	n := min(len(x), len(y))
	if n == 0 {
		return 0
	}
	if n > len(dst) {
		panic("subtle.XORBytes: dst too short")
	}
	if alias.InexactOverlap(dst[:n], x[:n]) || alias.InexactOverlap(dst[:n], y[:n]) {
		panic("subtle.XORBytes: invalid overlap")
	}
	dst, x, y = dst[:n], x[:n], y[:n]
	for i := range dst {
		dst[i] = x[i] ^ y[i]
	}
	return n
}
