package building

import "image/color"

// Smithing are what a smith holds at the anvil, laid on the ground as the
// items are, long along X with the working end toward +X: the hammer, a bar
// of iron, tongs, and tongs holding the bar.
var Smithing = []Building{Hammer, IronBar, Tongs, TongsIron}

// The colours of hot iron: dull orange where it has begun to cool, and
// bright where it is hottest.
var (
	hotIron  = paint{of: Flat, tint: color.NRGBA{232, 96, 24, 255}}
	whiteHot = paint{of: Flat, tint: color.NRGBA{255, 164, 52, 255}}
)

// Hammer is a smith's cross-peen hammer: a handle of 35 cm, its iron head
// of about a kilogram across its end, the flat face toward -Z, the peen
// that tapers to an edge toward +Z.
var Hammer = Building{Name: "hammer", Width: 0.4, Depth: 0.14, build: func(b *builder) {
	wood, iron := paint{of: Timber}, paint{of: Iron}
	b.lay(beam{a: vec{-0.2, 0.03, 0}, b: vec{0.15, 0.03, 0}, w: 0.03, h: 0.03, paint: wood})
	b.block(box{lo: vec{0.135, 0.012, -0.05}, hi: vec{0.165, 0.05, 0.02}, paint: iron})
	b.convex(hexa([8]vec{
		{0.135, 0.012, 0.02}, {0.165, 0.012, 0.02}, {0.155, 0.012, 0.07}, {0.145, 0.012, 0.07},
		{0.135, 0.05, 0.02}, {0.165, 0.05, 0.02}, {0.155, 0.05, 0.07}, {0.145, 0.05, 0.07},
	}, iron))
}}

// IronBar is a bar of iron 30 cm long, square, cold at one end and
// glowing at the other.
var IronBar = Building{Name: "iron_bar", Width: 0.32, Depth: 0.03, build: func(b *builder) {
	b.bar(-0.15, 0.15)
}}

// bar adds a bar of iron from x0 to x1 along X, its half toward +X hot.
func (b *builder) bar(x0, x1 float64) {
	mid := (x0 + x1) / 2
	b.block(box{lo: vec{x0, 0, -0.0125}, hi: vec{mid, 0.025, 0.0125}, paint: paint{of: Iron}})
	b.block(box{lo: vec{mid, 0, -0.0125}, hi: vec{(mid + x1) / 2, 0.025, 0.0125}, paint: hotIron})
	b.block(box{lo: vec{(mid + x1) / 2, 0, -0.0125}, hi: vec{x1, 0.025, 0.0125}, paint: whiteHot})
}

// Tongs are a smith's tongs: two reins of 30 cm that meet at a rivet and
// go on as jaws, closed on what they hold.
var Tongs = Building{Name: "tongs", Width: 0.42, Depth: 0.12, build: func(b *builder) { b.tongs(0) }}

// TongsIron are tongs holding a bar of iron by its cold end, the hot end
// out beyond the jaws.
var TongsIron = Building{Name: "tongs_iron", Width: 0.62, Depth: 0.12, build: func(b *builder) {
	b.tongs(0)
	b.bar(0.0, 0.28)
}}

// tongs adds tongs whose rivet is at x along X, its reins behind and its
// jaws ahead.
func (b *builder) tongs(x float64) {
	iron := paint{of: Iron}
	for _, side := range []float64{-1, 1} {
		b.lay(beam{a: vec{x - 0.3, 0.0125, side * 0.05}, b: vec{x, 0.0125, side * 0.013}, w: 0.016, h: 0.02, paint: iron})
		b.lay(beam{a: vec{x, 0.0125, side * 0.013}, b: vec{x + 0.1, 0.0125, side * 0.0125}, w: 0.012, h: 0.02, paint: iron})
	}
	b.turn(column{a: vec{x, 0, 0}, b: vec{x, 0.026, 0}, ra: 0.01, rb: 0.01, sides: 8, paint: iron})
}
