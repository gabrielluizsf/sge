//go:build gots_overlay

// Overlay for the P-256 generator tables, which Go aliases into the
// embedded bytes with unsafe: here they are decoded from them instead,
// each coordinate fiat's Montgomery form, little-endian, x then y for
// every point.

package nistec

import "crypto/internal/fips140/nistec/fiat"

// init is the package's only one, which builds p256GeneratorTables.
func init() {
	t := new([43]p256AffineTable)
	b := p256PrecomputedEmbed[:]
	for i := range t {
		for j := range t[i] {
			t[i][j].x.GotsSetMontgomery(b[:32])
			t[i][j].y.GotsSetMontgomery(b[32:64])
			b = b[64:]
		}
	}
	p256GeneratorTables = t
}

var _ = fiat.P256Element{}
