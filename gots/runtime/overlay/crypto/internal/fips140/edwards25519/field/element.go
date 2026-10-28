//go:build gots_overlay

// Overlay for Element under gots. Go's works on 51-bit limbs in uint64s, a
// BigInt each here; this one is twelve 22-bit limbs in int32s, in the
// Montgomery domain, R = 2^264, always reduced below p, its arithmetic
// fe_arith.go's (see nistec/fiat/mont.go). Invert, Pow22523, SqrtRatio,
// IsNegative and Absolute are Go's, on these methods.

package field

// Element represents an element of the field GF(2^255-19). Note that this
// is not a cryptographically secure group, and should only be used to interact
// with edwards25519.Point coordinates.
//
// This type works similarly to math/big.Int, and all arguments and receivers
// are allowed to alias.
//
// The zero value is a valid zero element.
type Element struct {
	l0, l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11 int32
}

const (
	limbBits = 22
	mask     = 1<<limbBits - 1

	// See nistec/fiat/mont.go.
	inv   = 1.0 / (1 << limbBits)
	base  = 1 << limbBits
	round = 3 << 73
	half  = base/2 - 0.5
)

var feZero = &Element{}

var feOne = &feR

// Set sets v = a, and returns v.
func (v *Element) Set(a *Element) *Element {
	*v = *a
	return v
}

// Equal returns 1 if v and u are equal, and 0 otherwise.
func (v *Element) Equal(u *Element) int {
	acc := (v.l0 ^ u.l0) | (v.l1 ^ u.l1) | (v.l2 ^ u.l2) | (v.l3 ^ u.l3) | (v.l4 ^ u.l4) | (v.l5 ^ u.l5) |
		(v.l6 ^ u.l6) | (v.l7 ^ u.l7) | (v.l8 ^ u.l8) | (v.l9 ^ u.l9) | (v.l10 ^ u.l10) | (v.l11 ^ u.l11)
	return 1 - int(uint32(acc|-acc)>>31)
}
