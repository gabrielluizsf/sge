//go:build gots_overlay

// Overlay for the inner digest crypto/internal/fips140hash reaches by
// linkname: exported for it, as transpiled packages cannot link names.

package sha3

import "crypto/internal/fips140/sha3"

// GotsUnwrap is the FIPS 140 digest inside s.
func GotsUnwrap(s *SHA3) *sha3.Digest {
	s.init()
	return &s.s
}
