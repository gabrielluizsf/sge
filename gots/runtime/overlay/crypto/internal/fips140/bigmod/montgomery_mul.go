//go:build gots_overlay

// Overlay for montgomeryMul, Go's algorithm for every size, its loop done
// by the runtime (rt/bits.ts) on the words' arrays as they are, in a
// scratch array kept from one call to the next: under gots each of Go's
// subslices and stack arrays would be a new object, and making them takes
// longer than the multiplication of P-256's scalars takes to compute. A
// call never waits, so one at a time uses the scratch.

package bigmod

var montgomeryScratch montgomeryJob

func (x *Nat) montgomeryMul(a *Nat, b *Nat, m *Modulus) *Nat {
	n := len(m.nat.limbs)
	j := &montgomeryScratch
	if cap(j.T) < n*2 {
		j.T = make([]uint, n*2)
	}
	j.T = j.T[:n*2]
	clear(j.T)
	j.a, j.b, j.m, j.m0inv = a.limbs[:n], b.limbs[:n], m.nat.limbs[:n], m.m0inv
	c := montgomeryWords(j)
	copy(x.reset(n).limbs, j.T[n:])
	x.maybeSubtractModulus(choice(c), m)
	return x
}
