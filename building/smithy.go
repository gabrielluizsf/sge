package building

import (
	"image/color"
	"math"
)

// Smithy is a village smithy, 5 by 4 m: an open shed of timber under a
// tiled gable, boarded at its back and left; a stone forge in its back
// corner, its hood narrowing into a chimney through the roof, coals
// glowing on its hearth; a quenching trough, the smith's tools on the
// wall, and the anvil out in front, where the light is.
var Smithy = Building{Name: "smithy", Width: 5, Depth: 4, build: smithy}

// The smithy's shed: half its width and depth, and its eaves' height.
const smithyW, smithyD, smithyEaves = 2.5, 2, 2.7

func smithy(b *builder) {
	shed(b)
	forge(b)
	trough(b, vec{0.1, 0, -1.55})
	for i, l := range []float64{0.55, 0.42, 0.62} {
		x := 0.9 + float64(i)*0.22
		b.lay(beam{a: vec{x, 1.75, -smithyD + 0.14}, b: vec{x, 1.75 - l, -smithyD + 0.14}, w: 0.03, h: 0.03, paint: paint{of: Iron}})
	}
	b.anvil(vec{0.7, 0, 1.3})
}

// shed is the smithy's frame, its boarded walls, and its roof.
func shed(b *builder) {
	wood := paint{of: Timber}
	w, d, e := smithyW-0.1, smithyD-0.1, smithyEaves
	for _, x := range []float64{-w, w} {
		for _, z := range []float64{-d, d} {
			b.block(box{lo: vec{x - 0.1, 0, z - 0.1}, hi: vec{x + 0.1, e, z + 0.1}, paint: wood})
		}
		b.lay(beam{a: vec{x, e - 0.1, -d - 0.1}, b: vec{x, e - 0.1, d + 0.1}, w: 0.18, h: 0.2, paint: wood})
	}
	for _, z := range []float64{-d, d} {
		b.lay(beam{a: vec{-w - 0.1, e + 0.1, z}, b: vec{w + 0.1, e + 0.1, z}, w: 0.2, h: 0.2, paint: wood})
	}
	boards := paint{of: Planks}
	b.block(box{lo: vec{-w, 0, -d - 0.06}, hi: vec{w, e, -d + 0.02}, paint: boards})
	b.block(box{lo: vec{-w - 0.06, 0, -d}, hi: vec{-w + 0.02, e, d}, paint: boards})
	p := pitched{w: 2 * smithyW, d: 2 * smithyD, eaves: e + 0.2, slope: math.Tan(32 * math.Pi / 180), over: 0.35}
	b.gable(cover{p, paint{of: Tiles}})
	b.gableEnds(cover{p, boards})
}

// forge is the smith's hearth of stone in the back left corner, its
// coals glowing, and over it a hood narrowing into a chimney up through
// the roof.
func forge(b *builder) {
	stone := paint{of: Stone}
	x0, x1, z0, z1 := -smithyW+0.1, -smithyW+1.5, -smithyD+0.1, -smithyD+1.15
	b.block(box{lo: vec{x0, 0, z0}, hi: vec{x1, 0.82, z1}, paint: stone})
	b.block(box{lo: vec{x0 + 0.2, 0.82, z0 + 0.15}, hi: vec{x1 - 0.2, 0.86, z1 - 0.2}, paint: paint{of: Soot}})
	glow := paint{of: Flat, tint: color.NRGBA{255, 112, 30, 255}}
	for i, at := range [][2]float64{{0.45, 0.4}, {0.62, 0.55}, {0.8, 0.42}, {0.55, 0.7}, {0.9, 0.62}} {
		r := 0.05 + 0.012*float64(i%2)
		b.roll(ball{c: vec{x0 + at[0], 0.86, z0 + at[1]}, r: r, squash: 0.5, paint: glow})
	}
	cx, cz := (x0+x1)/2, (z0+z1)/2
	b.raiseTaper(taper{at: vec{cx, 1.55, cz}, h: 0.9, base: [2]float64{(x1 - x0) / 2, (z1 - z0) / 2}, top: [2]float64{0.3, 0.28}, paint: stone})
	b.block(box{lo: vec{cx - 0.3, 2.45, cz - 0.28}, hi: vec{cx + 0.3, 4.75, cz + 0.28}, paint: stone})
	b.block(box{lo: vec{cx - 0.2, 4.6, cz - 0.18}, hi: vec{cx + 0.2, 4.755, cz + 0.18}, paint: paint{of: Soot}})
	for _, x := range []float64{x0 + 0.05, x1 - 0.05} {
		b.block(box{lo: vec{x - 0.05, 0.82, z0}, hi: vec{x + 0.05, 1.55, z0 + 0.12}, paint: stone})
	}
}

// trough is a quenching trough of boards, full of water, its foot's
// middle at at.
func trough(b *builder, at vec) {
	boards := paint{of: Planks}
	o := func(x, y, z float64) vec { return add(at, vec{x, y, z}) }
	const w, d, h, t = 0.55, 0.25, 0.6, 0.05
	b.block(box{lo: o(-w, 0, -d), hi: o(w, t, d), paint: boards})
	b.block(box{lo: o(-w, 0, -d), hi: o(w, h, -d+t), paint: boards})
	b.block(box{lo: o(-w, 0, d-t), hi: o(w, h, d), paint: boards})
	b.block(box{lo: o(-w, 0, -d+t), hi: o(-w+t, h, d-t), paint: boards})
	b.block(box{lo: o(w-t, 0, -d+t), hi: o(w, h, d-t), paint: boards})
	b.block(box{lo: o(-w+t, t, -d+t), hi: o(w-t, h-0.08, d-t), paint: paint{of: Water}})
}
