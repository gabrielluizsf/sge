//go:build gots_overlay

// Overlay for the sha3 digest crypto/sha3 hands over by linkname.

package fips140hash

import (
	"crypto/internal/fips140/sha3"
	sha3pub "crypto/sha3"
)

func sha3Unwrap(s *sha3pub.SHA3) *sha3.Digest { return sha3pub.GotsUnwrap(s) }
