package building

import "image/color"

// Stall is a market stall, 2.6 by 1.8 m: a boarded counter at its
// front, crates of oranges, apples and tomatoes on it, and a striped
// canvas awning on four posts, sloping toward the buyers.
var Stall = Building{Name: "stall", Width: 2.6, Depth: 1.8, build: stall}

// produce is the colours of what a stall sells, in sRGB.
var produce = []color.NRGBA{{232, 132, 28, 255}, {150, 180, 60, 255}, {196, 40, 30, 255}}

func stall(b *builder) {
	wood, boards := paint{of: Timber}, paint{of: Planks}
	const x, front, back = 1.2, 0.8, -0.8
	for _, post := range []vec{{-x, 2.25, front}, {x, 2.25, front}, {-x, 2.6, back}, {x, 2.6, back}} {
		b.block(box{lo: vec{post[0] - 0.05, 0, post[2] - 0.05}, hi: vec{post[0] + 0.05, post[1], post[2] + 0.05}, paint: wood})
	}
	b.lay(beam{a: vec{-x - 0.05, 2.2, front}, b: vec{x + 0.05, 2.2, front}, w: 0.08, h: 0.1, paint: wood})
	b.lay(beam{a: vec{-x - 0.05, 2.55, back}, b: vec{x + 0.05, 2.55, back}, w: 0.08, h: 0.1, paint: wood})
	// The counter: its top, its front and ends boarded.
	b.block(box{lo: vec{-x - 0.05, 0.88, 0.1}, hi: vec{x + 0.05, 0.93, front + 0.08}, paint: boards})
	b.block(box{lo: vec{-x, 0.05, front - 0.03}, hi: vec{x, 0.88, front + 0.02}, paint: boards})
	for _, sx := range []float64{-x, x - 0.05} {
		b.block(box{lo: vec{sx, 0.05, 0.12}, hi: vec{sx + 0.05, 0.88, front - 0.03}, paint: boards})
	}
	for i, c := range produce {
		crate(b, crateOf{at: vec{-0.75 + 0.75*float64(i), 0.93, 0.45}, fruit: paint{of: Flat, tint: c}})
	}
	canvas := paint{of: Canvas}
	b.convex(hexa([8]vec{
		{-x - 0.2, 2.6, back - 0.25}, {x + 0.2, 2.6, back - 0.25}, {x + 0.2, 2.2, front + 0.35}, {-x - 0.2, 2.2, front + 0.35},
		{-x - 0.2, 2.63, back - 0.25}, {x + 0.2, 2.63, back - 0.25}, {x + 0.2, 2.23, front + 0.35}, {-x - 0.2, 2.23, front + 0.35},
	}, canvas))
	b.block(box{lo: vec{-x - 0.2, 1.98, front + 0.33}, hi: vec{x + 0.2, 2.21, front + 0.36}, paint: canvas})
}

// crateOf is a crate of fruit: its foot's middle, and the fruit's paint.
type crateOf struct {
	at    vec
	fruit paint
}

// crate is a shallow crate of boards heaped with fruit.
func crate(b *builder, c crateOf) {
	o := func(x, y, z float64) vec { return add(c.at, vec{x, y, z}) }
	b.block(box{lo: o(-0.3, 0, -0.2), hi: o(0.3, 0.14, 0.2), paint: paint{of: Planks}})
	for i := range 5 {
		for k := range 4 {
			x, z := -0.24+0.12*float64(i), -0.15+0.1*float64(k)
			y := 0.14 + 0.02*float64((i+k)%2)
			b.roll(ball{c: o(x, y, z), r: 0.05, squash: 0.9, paint: c.fruit})
		}
	}
}
