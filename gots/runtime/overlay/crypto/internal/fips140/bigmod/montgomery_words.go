//go:build gots_overlay

// The loop of montgomeryMul: see Go's for how it works.

package bigmod

import "math/bits"

// montgomeryJob is a montgomeryMul's operands, T its 2n words, zero.
type montgomeryJob struct {
	T, a, b, m []uint
	m0inv      uint
}

// montgomeryWords leaves a·b/R in T's upper half, to be reduced by m once
// more if the carry it gives back is 1. The runtime's does it under gots:
// this is the same, for the type checker.
func montgomeryWords(j *montgomeryJob) uint {
	n := len(j.m)
	var c uint
	for i := 0; i < n; i++ {
		c1 := addMulVVW(j.T[i:n+i], j.a, j.b[i])
		Y := j.T[i] * j.m0inv
		c2 := addMulVVW(j.T[i:n+i], j.m, Y)
		j.T[n+i], c = bits.Add(c1, c2, c)
	}
	return c
}
