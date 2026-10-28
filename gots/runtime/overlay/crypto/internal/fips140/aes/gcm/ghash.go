//go:build gots_overlay

// Overlay for GHASH's arithmetic, Go's own (BearSSL's constant-time
// carry-less multiplication, Karatsuba, and the field's reduction) on
// 32-bit words only: Go's 64-bit words are each a BigInt under gots, which
// made AES-GCM seal some 0.4 MB a second.

package gcm

import "crypto/internal/fips140deps/byteorder"

// ghashKey is H, least significant word first.
type ghashKey struct{ h0, h1, h2, h3 uint32 }

// ghashState is Y, least significant word first.
type ghashState struct{ w0, w1, w2, w3 uint32 }

func ghash(out, H *[gcmBlockSize]byte, inputs ...[]byte) {
	clear(out[:])
	k := ghashKey{byteorder.BEUint32(H[12:]), byteorder.BEUint32(H[8:]), byteorder.BEUint32(H[4:]), byteorder.BEUint32(H[:4])}
	var y ghashState
	for _, input := range inputs {
		for len(input) >= 16 {
			ghashAbsorb(&y, &k, input[:16])
			input = input[16:]
		}
		if len(input) > 0 {
			var partial [gcmBlockSize]byte
			copy(partial[:], input)
			ghashAbsorb(&y, &k, partial[:])
		}
	}
	byteorder.BEPutUint32(out[0:], y.w3)
	byteorder.BEPutUint32(out[4:], y.w2)
	byteorder.BEPutUint32(out[8:], y.w1)
	byteorder.BEPutUint32(out[12:], y.w0)
}

// ghashAbsorb sets Y = (Y ^ b)·H: Go's ghash loop, step for step. The
// runtime's (rt/ghash.ts) does it under gots, on numbers, making no
// object: this is the same, for the type checker.
func ghashAbsorb(y *ghashState, k *ghashKey, b []byte) {
	y.w3 ^= byteorder.BEUint32(b[0:])
	y.w2 ^= byteorder.BEUint32(b[4:])
	y.w1 ^= byteorder.BEUint32(b[8:])
	y.w0 ^= byteorder.BEUint32(b[12:])
	r00, r01 := halves(ghashMul(y.w0, k.h0), ghashMul(y.w1, k.h1), ghashMul(y.w0^y.w1, k.h0^k.h1))
	r10, r11 := halves(ghashMul(y.w2, k.h2), ghashMul(y.w3, k.h3), ghashMul(y.w2^y.w3, k.h2^k.h3))
	r20, r21 := halves(ghashMul(y.w0^y.w2, k.h0^k.h2), ghashMul(y.w1^y.w3, k.h1^k.h3),
		ghashMul(y.w0^y.w2^y.w1^y.w3, k.h0^k.h2^k.h1^k.h3))
	r20 = r20.xor(r00).xor(r10)
	r21 = r21.xor(r01).xor(r11)
	r01 = r01.xor(r20)
	r10 = r10.xor(r21)
	// The 256-bit product, shifted by one bit, as GHASH has it.
	z0 := r00.shl(1)
	z1 := r01.shl(1).xor(r00.shr(63))
	z2 := r10.shl(1).xor(r01.shr(63))
	z3 := r11.shl(1).xor(r10.shr(63))
	// Reduced modulo the field polynomial.
	low, carry := fold(z0)
	z2, z1 = z2.xor(low), z1.xor(carry)
	low, carry = fold(z1)
	z3, z2 = z3.xor(low), z2.xor(carry)
	y.w0, y.w1, y.w2, y.w3 = z2.l, z2.h, z3.l, z3.h
}

// halves is the 128-bit product of Karatsuba's three, low and high.
func halves(lo, hi, sum u64) (u64, u64) {
	mid := sum.xor(lo).xor(hi)
	return lo.xor(mid.shl(32)), hi.xor(mid.shr(32))
}

// fold is what the high-degree word lw adds to the low-degree terms, and
// to the next word down.
func fold(lw u64) (low, carry u64) {
	return lw.xor(lw.shr(1)).xor(lw.shr(2)).xor(lw.shr(7)), lw.shl(63).xor(lw.shl(62)).xor(lw.shl(57))
}
