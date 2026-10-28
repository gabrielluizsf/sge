//go:build gots_overlay

// Element arithmetic, on fe_arith.go's: see element.go.

package field

// Add sets v = a + b, and returns v.
func (v *Element) Add(a, b *Element) *Element {
	feAdd(v, a, b)
	return v
}

// Subtract sets v = a - b, and returns v.
func (v *Element) Subtract(a, b *Element) *Element {
	feSub(v, a, b)
	return v
}

// Negate sets v = -a, and returns v.
func (v *Element) Negate(a *Element) *Element {
	feSub(v, feZero, a)
	return v
}

// Multiply sets v = x * y, and returns v.
func (v *Element) Multiply(x, y *Element) *Element {
	feMul(v, x, y)
	return v
}

// Square sets v = x * x, and returns v.
func (v *Element) Square(x *Element) *Element {
	feMul(v, x, x)
	return v
}

// SquareN sets v = x^(2^n), and returns v.
func (v *Element) SquareN(x *Element, n int) *Element {
	feMul(v, x, x)
	for range n - 1 {
		feMul(v, v, v)
	}
	return v
}

// Mult32 sets v = x * y, and returns v.
func (v *Element) Mult32(x *Element, y uint32) *Element {
	var t Element
	feMul(&t, &Element{l0: int32(y & mask), l1: int32(y >> limbBits)}, &feRR)
	feMul(v, x, &t)
	return v
}
