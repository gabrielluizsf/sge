//go:build gots_overlay

// The fields' arithmetic under gots. fiat's works on 64-bit limbs, each a
// BigInt here, and makes a P-256 signature check take a quarter of a
// second. Here an element is 22-bit limbs held in float64s, which
// JavaScript computes on as they are: a limb times a limb is 44 bits, and
// a float64 holds 53 exactly. An element keeps its limbs in fields,
// as a typed array takes some 40 times as long to make as an object, and
// works on them in the field's scratch arrays. Elements are in the
// Montgomery domain, R = 2^(22n), and always fully reduced, below p. There
// are no branches on the values, as there are none in fiat's.

package fiat

import "math"

const (
	maxLimbs = 24 // P-521's 521 bits
	limbBits = 22
	base     = 1 << limbBits
	mask     = base - 1

	// x - half + round - round is x rounded down to a multiple of base,
	// for x from 0 to 2^52 a whole number: the float64s from round - 2^52
	// to round + 2^52 are base apart, and x - half, never halfway between
	// two, as it ends in .5, rounds to the nearest.
	round = 3 << 73
	half  = base/2 - 0.5
	inv   = 1.0 / base
)

type limbs = [maxLimbs]float64

// wide holds a sum before its reduction: n limbs and two of carry.
type wide = [maxLimbs + 2]float64

type montField struct {
	n       int     // limbs in use
	p       limbs   // the modulus
	pInv    float64 // -p⁻¹ mod 2²²
	one     limbs   // R mod p, 1 in the Montgomery domain
	rr      limbs   // R² mod p, which takes an element into the domain
	byteLen int

	// Scratch for the arithmetic, which never waits, so one operation at a
	// time uses it: a, b and c for the elements' methods, the rest for mul,
	// reduce and sub.
	a, b, c limbs
	m, s, d limbs
	t       wide
}

// mul sets z = x*y/R mod p, by Montgomery multiplication scanning the
// product a column at a time: a column's sum of up to 2n products of 44
// bits fits a float64 exactly, so it takes one carry, not one a product.
func (f *montField) mul(z, x, y *limbs) {
	n := f.n
	p, m, t := &f.p, &f.m, &f.t
	c := 0.0
	for i := 0; i < n; i++ {
		acc := c + x[i]*y[0]
		for j := 0; j < i; j++ {
			acc += x[j]*y[i-j] + m[j]*p[i-j]
		}
		// m[i] makes the column's lowest limb zero.
		v := (acc - math.Floor(acc*inv)*base) * f.pInv
		m[i] = v - math.Floor(v*inv)*base
		c = (acc + m[i]*p[0]) * inv
	}
	for i := n; i < 2*n; i++ {
		acc := c
		for j := i - n + 1; j < n; j++ {
			acc += x[j]*y[i-j] + m[j]*p[i-j]
		}
		c = math.Floor(acc * inv)
		t[i-n] = acc - c*base
	}
	t[n] = c
	f.reduce(z, t)
}

// reduce sets z = t mod p, for t below 2p.
func (f *montField) reduce(z *limbs, t *wide) {
	n := f.n
	s := &f.s
	b := 0.0
	for j := 0; j < n; j++ {
		d := t[j] - f.p[j] - b
		b = -math.Floor(d * inv)
		s[j] = d + b*base
	}
	// b is 1 when t is below p, and t is kept.
	b = -math.Floor((t[n] - b) * inv)
	for j := 0; j < n; j++ {
		z[j] = s[j] + b*(t[j]-s[j])
	}
}

func (f *montField) add(z, x, y *limbs) {
	n := f.n
	t := &f.t
	c := 0.0
	for j := 0; j < n; j++ {
		v := x[j] + y[j] + c
		c = math.Floor(v * inv)
		t[j] = v - c*base
	}
	t[n] = c
	f.reduce(z, t)
}

func (f *montField) sub(z, x, y *limbs) {
	n := f.n
	d := &f.d
	b := 0.0
	for j := 0; j < n; j++ {
		v := x[j] - y[j] - b
		b = -math.Floor(v * inv)
		d[j] = v + b*base
	}
	// Below zero: add p back.
	c := 0.0
	for j := 0; j < n; j++ {
		v := d[j] + b*f.p[j] + c
		c = math.Floor(v * inv)
		z[j] = v - c*base
	}
}

// selectLimbs sets z to a if cond is 1, and to b if it is 0.
func selectLimbs(z, a, b *limbs, cond int) {
	k := float64(cond)
	for j := range z {
		z[j] = b[j] + k*(a[j]-b[j])
	}
}

// equal is 1 if x and y, both reduced, are the same, and 0 otherwise.
func equal(x, y *limbs) int {
	acc := 0.0
	for j := range x {
		d := x[j] - y[j]
		acc += d * d
	}
	return 1 - int(math.Ceil(acc/(acc+1)))
}
