package tree

import "math"

// ball is a round fruit, a coconut or a cashew apple: its middle, its
// radius across, how tall it is as a share of that, and its material.
type ball struct {
	at     vec
	r      float64
	squash float64
	mat    int
}

// ball adds b.
func (g *grower) ball(b ball) {
	const rings, segs = 8, 10
	p := g.prim(b.mat)
	first := uint32(len(p.Positions))
	for r := 0; r <= rings; r++ {
		th := math.Pi * float64(r) / rings
		for s := 0; s <= segs; s++ {
			ph := 2 * math.Pi * float64(s) / segs
			n := vec{math.Sin(th) * math.Cos(ph), math.Cos(th), math.Sin(th) * math.Sin(ph)}
			p.Positions = append(p.Positions, f32(add(b.at, vec{n[0] * b.r, n[1] * b.r * b.squash, n[2] * b.r})))
			p.Normals = append(p.Normals, f32(n))
			p.UVs = append(p.UVs, [2]float32{float32(s) / segs, float32(r) / rings})
			if r > 0 && s > 0 {
				a := first + uint32((r-1)*(segs+1)+s-1)
				c := a + uint32(segs+1)
				p.Indices = append(p.Indices, a, c, c+1, a, c+1, a+1)
			}
		}
	}
}
