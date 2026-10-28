//go:build gots_overlay

// GHASH's carry-less multiplication on 32-bit words: see ghash.go.

package gcm

import "math/bits"

// spread is a word in four, each with every fourth bit of it.
type spread struct{ m0, m1, m2, m3 uint32 }

func spreadOf(x uint32) spread {
	return spread{x & 0x11111111, x & 0x22222222, x & 0x44444444, x & 0x88888888}
}

func (s *spread) at(i int) uint32 {
	switch i & 3 {
	case 0:
		return s.m0
	case 1:
		return s.m1
	case 2:
		return s.m2
	}
	return s.m3
}

// term is the sum of the products landing on bits i mod 4: x times the
// circulant matrix of y's row i, whose carries are masked off after.
func (x *spread) term(y *spread, i int) u64 {
	h0, l0 := bits.Mul32(x.m0, y.at(i))
	h1, l1 := bits.Mul32(x.m1, y.at(i+3))
	h2, l2 := bits.Mul32(x.m2, y.at(i+2))
	h3, l3 := bits.Mul32(x.m3, y.at(i+1))
	return u64{h0 ^ h1 ^ h2 ^ h3, l0 ^ l1 ^ l2 ^ l3}
}

// ghashMul is x·y carry-less: Go's, from integer products of x and y with
// every fourth bit kept, whose carries fall in the bits between.
func ghashMul(x, y uint32) u64 {
	xs, ys := spreadOf(x), spreadOf(y)
	var z u64
	for i, m := range [4]uint32{0x11111111, 0x22222222, 0x44444444, 0x88888888} {
		t := xs.term(&ys, i)
		z.h |= t.h & m
		z.l |= t.l & m
	}
	return z
}
