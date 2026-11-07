package building

import "image/color"

// Bakery is a town's bakery, 5 by 4 m: the open shed, a domed brick oven at
// its back with the fire in its mouth, and a counter of loaves at the front.
var Bakery = Building{Name: "bakery", Width: 5, Depth: 4, build: func(b *builder) {
	shed(b)
	brick := paint{of: Brick}
	b.turn(column{a: vec{-0.8, 0, -1.1}, b: vec{-0.8, 0.9, -1.1}, ra: 0.75, rb: 0.72, sides: 16, paint: brick})
	b.roll(ball{c: vec{-0.8, 0.9, -1.1}, r: 0.72, squash: 0.75, paint: brick})
	b.block(box{lo: vec{-1.05, 0.3, -0.42}, hi: vec{-0.55, 0.75, -0.36}, paint: paint{of: Flat, tint: color.NRGBA{240, 110, 30, 255}}})
	frontCounter(b)
	loaf := paint{of: Flat, tint: color.NRGBA{200, 145, 70, 255}}
	for i := range 5 {
		b.roll(ball{c: vec{-1.4 + 0.42*float64(i), 1.06, 1.52}, r: 0.15, squash: 0.65, paint: loaf})
	}
}}

// Butcher is a town's butcher's, 5 by 4 m: the open shed, sides of meat
// hanging from a rail at its back, a block and a counter of cuts at the front.
var Butcher = Building{Name: "butcher", Width: 5, Depth: 4, build: func(b *builder) {
	shed(b)
	red := paint{of: Flat, tint: color.NRGBA{140, 44, 38, 255}}
	b.lay(beam{a: vec{-1.8, 2.2, -1.2}, b: vec{1.8, 2.2, -1.2}, w: 0.05, h: 0.05, paint: paint{of: Iron}})
	for i := range 4 {
		x := -1.3 + 0.9*float64(i)
		b.lay(beam{a: vec{x, 2.2, -1.2}, b: vec{x, 2.0, -1.2}, w: 0.02, h: 0.02, paint: paint{of: Iron}})
		b.roll(ball{c: vec{x, 1.55, -1.2}, r: 0.26, squash: 1.7, paint: red})
	}
	b.turn(column{a: vec{1.3, 0, 0.6}, b: vec{1.3, 0.85, 0.6}, ra: 0.36, rb: 0.36, sides: 14, paint: paint{of: Timber}})
	frontCounter(b)
	for i := range 4 {
		b.block(box{lo: vec{-1.5 + 0.5*float64(i), 1.01, 1.45}, hi: vec{-1.15 + 0.5*float64(i), 1.07, 1.62}, paint: red})
	}
}}

// frontCounter is the boarded counter at a shed's front, where one buys.
func frontCounter(b *builder) {
	b.block(box{lo: vec{-1.6, 0, 1.45}, hi: vec{0.4, 0.95, 1.6}, paint: paint{of: Planks}})
	b.block(box{lo: vec{-1.7, 0.95, 1.35}, hi: vec{0.5, 1.01, 1.7}, paint: paint{of: Timber}})
}

// What the animals give and the bakeries make, as items carried.
var (
	// WoolBall is a ball of wool.
	WoolBall = Building{Name: "wool", Width: 0.2, Depth: 0.2, build: func(b *builder) {
		b.roll(ball{c: vec{0, 0.1, 0}, r: 0.1, squash: 1, paint: paint{of: Flat, tint: color.NRGBA{238, 232, 216, 255}}})
	}}
	// Loaf is a loaf of bread.
	Loaf = Building{Name: "bread", Width: 0.3, Depth: 0.15, build: func(b *builder) {
		b.roll(ball{c: vec{0, 0.085, 0}, r: 0.14, squash: 0.6, paint: paint{of: Flat, tint: color.NRGBA{200, 145, 70, 255}}})
	}}
	// BeefCut is a cut of beef.
	BeefCut = Building{Name: "beef", Width: 0.26, Depth: 0.2, build: func(b *builder) {
		b.block(box{lo: vec{-0.13, 0, -0.1}, hi: vec{0.13, 0.08, 0.1}, paint: paint{of: Flat, tint: color.NRGBA{158, 46, 42, 255}}})
	}}
	// MuttonCut is a cut of mutton.
	MuttonCut = Building{Name: "mutton", Width: 0.22, Depth: 0.18, build: func(b *builder) {
		b.block(box{lo: vec{-0.11, 0, -0.09}, hi: vec{0.11, 0.08, 0.09}, paint: paint{of: Flat, tint: color.NRGBA{190, 90, 84, 255}}})
	}}
)
