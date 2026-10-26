//go:build gots_overlay

// Overlay for the P-256 generator tables, which Go aliases into the
// embedded bytes with unsafe: here they are decoded from them instead,
// each limb a little-endian uint64, x's four then y's for every point.

package nistec

import (
	"crypto/internal/fips140/nistec/fiat"
	"crypto/internal/fips140deps/byteorder"
)

// init is the package's only one, which builds p256GeneratorTables.
func init() {
	t := new([43]p256AffineTable)
	b := p256PrecomputedEmbed[:]
	for i := range t {
		for j := range t[i] {
			t[i][j].x.GotsSetLimbs(limbs(b))
			t[i][j].y.GotsSetLimbs(limbs(b[32:]))
			b = b[64:]
		}
	}
	p256GeneratorTables = t
}

func limbs(b []byte) [4]uint64 {
	return [4]uint64{byteorder.LEUint64(b), byteorder.LEUint64(b[8:]), byteorder.LEUint64(b[16:]), byteorder.LEUint64(b[24:])}
}

var _ = fiat.P256Element{}
