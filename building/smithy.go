package building

import (
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
	b.forge(hearth{x0: -smithyW + 0.1, x1: -smithyW + 1.5, z0: -smithyD + 0.1, z1: -smithyD + 1.15, top: 4.75})
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
