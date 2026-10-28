package main

// arith is the unrolled arithmetic of p's elements. Sums and differences
// of limbs are int32s, whose carries JavaScript takes with a shift; a
// product's column, up to 2n products of 44 bits, is a float64, summed as
// a tree, not a chain, as each sum would otherwise wait for the one before,
// and its carry taken by rounding (see mont.go's round), not by
// math.Floor, which takes JavaScript some 4 times as long.
func arith(p *prime) []byte {
	var o out
	o.header(p)
	mul(&o, p)
	add(&o, p)
	sub(&o, p)
	return o.Bytes()
}
