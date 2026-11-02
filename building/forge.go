package building

import "image/color"

// Forge is a smith's forge standing alone: a stone hearth, coals glowing
// on it, under a hood narrowing into a chimney, about a metre square.
var Forge = Building{Name: "forge", Width: 1, Depth: 0.9, build: func(b *builder) {
	b.forge(hearth{x0: -0.5, x1: 0.5, z0: -0.45, z1: 0.45, top: 3})
}}

// hearth is a forge's stone hearth from x0, z0 to x1, z1 on the ground,
// and how high the top of its chimney is.
type hearth struct{ x0, x1, z0, z1, top float64 }

// forge builds a forge on hearth h: its hearth, the coals on it, the hood
// over them on two posts, and its chimney.
func (b *builder) forge(h hearth) {
	stone := paint{of: Stone}
	x0, x1, z0, z1 := h.x0, h.x1, h.z0, h.z1
	b.block(box{lo: vec{x0, 0, z0}, hi: vec{x1, 0.82, z1}, paint: stone})
	b.block(box{lo: vec{x0 + 0.2, 0.82, z0 + 0.15}, hi: vec{x1 - 0.2, 0.86, z1 - 0.2}, paint: paint{of: Soot}})
	glow := paint{of: Flat, tint: color.NRGBA{255, 112, 30, 255}}
	for i, at := range [][2]float64{{0.45, 0.4}, {0.62, 0.55}, {0.8, 0.42}, {0.55, 0.7}, {0.9, 0.62}} {
		if x0+at[0] > x1-0.25 || z0+at[1] > z1-0.25 {
			continue // past a small hearth's bed of coals
		}
		r := 0.05 + 0.012*float64(i%2)
		b.roll(ball{c: vec{x0 + at[0], 0.86, z0 + at[1]}, r: r, squash: 0.5, paint: glow})
	}
	cx, cz := (x0+x1)/2, (z0+z1)/2
	b.raiseTaper(taper{at: vec{cx, 1.55, cz}, h: 0.9, base: [2]float64{(x1 - x0) / 2, (z1 - z0) / 2}, top: [2]float64{0.3, 0.28}, paint: stone})
	b.block(box{lo: vec{cx - 0.3, 2.45, cz - 0.28}, hi: vec{cx + 0.3, h.top, cz + 0.28}, paint: stone})
	b.block(box{lo: vec{cx - 0.2, h.top - 0.15, cz - 0.18}, hi: vec{cx + 0.2, h.top + 0.005, cz + 0.18}, paint: paint{of: Soot}})
	for _, x := range []float64{x0 + 0.05, x1 - 0.05} {
		b.block(box{lo: vec{x - 0.05, 0.82, z0}, hi: vec{x + 0.05, 1.55, z0 + 0.12}, paint: stone})
	}
}
