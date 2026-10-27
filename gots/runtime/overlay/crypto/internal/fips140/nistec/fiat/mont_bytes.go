//go:build gots_overlay

// Elements to and from their big-endian encodings: see mont.go.

package fiat

import "math"

// plain sets x to v's integer, big-endian, without reducing it. v fits
// the field's limbs.
func (f *montField) plain(x *limbs, v []byte) {
	*x = limbs{}
	acc, scale, j := 0.0, 1.0, 0
	for i := len(v) - 1; i >= 0; i-- {
		acc += float64(v[i]) * scale
		scale *= 256
		if scale >= base {
			hi := math.Floor(acc * inv)
			x[j] = acc - hi*base
			j++
			acc, scale = hi, scale*inv
		}
	}
	if j < f.n {
		x[j] = acc
	}
}

// below is 1 if x, f.n limbs, is less than p, and 0 otherwise.
func (f *montField) below(x *limbs) int {
	b := 0.0
	for j := 0; j < f.n; j++ {
		b = -math.Floor((x[j] - f.p[j] - b) * inv)
	}
	return int(b)
}

// setBytes sets z to v, a big-endian encoding of f.byteLen bytes, and is
// false, leaving z as it is, if v is not that long or not below p.
func (f *montField) setBytes(z *limbs, v []byte) bool {
	if len(v) != f.byteLen {
		return false
	}
	var x limbs
	f.plain(&x, v)
	if f.below(&x) == 0 {
		return false
	}
	f.mul(z, &x, &f.rr)
	return true
}

// bytes writes x's big-endian encoding to out, f.byteLen bytes.
func (f *montField) bytes(out []byte, x *limbs) {
	var v, one limbs
	one[0] = 1
	f.mul(&v, x, &one) // out of the Montgomery domain
	acc, bits, j := 0.0, 0, 0
	for i := len(out) - 1; i >= 0; i-- {
		if bits < 8 && j < f.n {
			acc += v[j] * float64(int(1)<<bits)
			bits += limbBits
			j++
		}
		hi := math.Floor(acc / 256)
		out[i] = byte(acc - hi*256)
		acc = hi
		bits -= 8
	}
}

// newMontField is the field of the prime p, big-endian, in n limbs.
func newMontField(n int, p []byte) *montField {
	f := &montField{n: n, byteLen: len(p)}
	f.plain(&f.p, p)
	// -p⁻¹ mod 2³², by Newton's iteration, each step doubling the bits
	// right, from the 3 an odd number is its own inverse to.
	p0 := uint32(f.p[0])
	r := p0
	for range 4 {
		r *= 2 - p0*r
	}
	f.pInv = float64(-r & (base - 1))
	// R mod p and R² mod p, by doubling 1.
	f.one[0] = 1
	for range limbBits * n {
		f.add(&f.one, &f.one, &f.one)
	}
	f.rr = f.one
	for range limbBits * n {
		f.add(&f.rr, &f.rr, &f.rr)
	}
	return f
}
