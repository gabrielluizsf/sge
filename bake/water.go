package bake

import "math"

// Water makes the ground's plane the surface of water, as a swimmer is
// baked: what lies under it is seen through it, tinted its Color, linear,
// and fading into it, gone a third of the way at Clear deep, in the
// model's units; it casts no shadow on the ground, which the game draws
// as water.
type Water struct {
	Color [3]float64
	Clear float64
}

// under is the sample c seen from above through depth of the water.
func (wt *Water) under(c [4]float64, depth float64) [4]float64 {
	k := math.Exp(-depth / max(wt.Clear, 1e-9))
	for i := range 3 {
		c[i] = c[i]*k + wt.Color[i]*(1-k)
	}
	c[3] *= k
	return c
}
