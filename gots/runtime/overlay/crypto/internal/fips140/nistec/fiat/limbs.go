//go:build gots_overlay

// Overlay for nistec's precomputed P-256 table, which Go reads by casting
// bytes to elements with unsafe: here each is decoded from them.

package fiat

// montRatio is mont.go's R over fiat's, 2^(22·12)/2^256, in the domain.
var montRatio = new(P256Element).setPower(limbBits*12 - 256)

// setPower sets e = 2^k, for k below 256.
func (e *P256Element) setPower(k int) *P256Element {
	var v [p256ElementLen]byte
	v[len(v)-1-k/8] = 1 << (k % 8)
	e.SetBytes(v[:])
	return e
}

// GotsSetMontgomery sets e to the element whose Montgomery form, as fiat
// keeps it, x·2^256 mod p, is the little-endian b, as the table stores it.
func (e *P256Element) GotsSetMontgomery(b []byte) *P256Element {
	var be [p256ElementLen]byte
	for i := range be {
		be[i] = b[len(be)-1-i]
	}
	f := p256Field
	f.plain(&f.a, be[:])
	montRatio.load(&f.b)
	f.mul(&f.c, &f.a, &f.b)
	e.store(&f.c)
	return e
}
