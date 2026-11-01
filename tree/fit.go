package tree

import "math"

// fit scales the tree, about its foot, to be exactly g.size tall, or
// wide where wide: the grown shape is true to its kind, and this sets
// its size true to the range.
func (g *grower) fit(wide bool) {
	lo, hi := vec{math.Inf(1), math.Inf(1), math.Inf(1)}, vec{math.Inf(-1), math.Inf(-1), math.Inf(-1)}
	for _, p := range g.prims {
		for _, q := range p.Positions {
			for k := range q {
				lo[k] = math.Min(lo[k], float64(q[k]))
				hi[k] = math.Max(hi[k], float64(q[k]))
			}
		}
	}
	got := hi[1]
	if wide {
		got = math.Max(hi[0]-lo[0], hi[2]-lo[2])
	}
	if got <= 0 {
		return
	}
	k := float32(g.size / got)
	for i := range g.prims {
		for j, q := range g.prims[i].Positions {
			g.prims[i].Positions[j] = [3]float32{q[0] * k, q[1] * k, q[2] * k}
		}
	}
}
