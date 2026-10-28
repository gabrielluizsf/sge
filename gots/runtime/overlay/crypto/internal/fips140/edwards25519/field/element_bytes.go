//go:build gots_overlay

// Elements to and from their encodings: see element.go.

package field

import "errors"

// SetBytes sets v to x, where x is a 32-byte little-endian encoding. If x is
// not of the right length, SetBytes returns nil and an error, and the
// receiver is unchanged.
//
// Consistent with RFC 7748, the most significant bit (the high bit of the
// last byte) is ignored, and non-canonical values (2^255-19 through 2^255-1)
// are accepted. Note that this is laxer than specified by RFC 8032, but
// consistent with most Ed25519 implementations.
func (v *Element) SetBytes(x []byte) (*Element, error) {
	if len(x) != 32 {
		return nil, errors.New("edwards25519: invalid field element input size")
	}
	var l [12]int32
	var acc uint32
	bits, j := 0, 0
	for i, b := range x {
		if i == 31 {
			b &= 0x7f
		}
		acc |= uint32(b) << bits
		bits += 8
		if bits >= limbBits {
			l[j] = int32(acc & mask)
			j++
			acc >>= limbBits
			bits -= limbBits
		}
	}
	l[j] = int32(acc)
	t := Element{l[0], l[1], l[2], l[3], l[4], l[5], l[6], l[7], l[8], l[9], l[10], l[11]}
	// Below 2^255, so 2p: the product is reduced.
	feMul(v, &t, &feRR)
	return v, nil
}

// Bytes returns the canonical 32-byte little-endian encoding of v.
func (v *Element) Bytes() []byte {
	var out [32]byte
	return v.bytes(&out)
}

func (v *Element) bytes(out *[32]byte) []byte {
	var t Element
	feMul(&t, v, &Element{l0: 1}) // out of the domain
	l := [12]int32{t.l0, t.l1, t.l2, t.l3, t.l4, t.l5, t.l6, t.l7, t.l8, t.l9, t.l10, t.l11}
	var acc uint32
	bits, j := 0, 0
	for i := range out {
		if bits < 8 && j < len(l) {
			acc |= uint32(l[j]) << bits
			bits += limbBits
			j++
		}
		out[i] = byte(acc)
		acc >>= 8
		bits -= 8
	}
	return out[:]
}

// sqrtM1 is 2^((p-1)/4), which squared is equal to -1 by Euler's Criterion.
var sqrtM1, _ = new(Element).SetBytes([]byte{
	0xb0, 0xa0, 0x0e, 0x4a, 0x27, 0x1b, 0xee, 0xc4, 0x78, 0xe4, 0x2f, 0xad, 0x06, 0x18, 0x43, 0x2f,
	0xa7, 0xd7, 0xfb, 0x3d, 0x99, 0x00, 0x4d, 0x2b, 0x0b, 0xdf, 0xc1, 0x4f, 0x80, 0x24, 0x83, 0x2b})
