package building

import "image/color"

// Stream is the water a can pours: a thin column of a metre standing on
// its origin, to be stretched from the ground to the spout, pale and
// bluish as water in the light.
var Stream = Building{Name: "stream", Width: 0.03, Depth: 0.03, build: func(b *builder) {
	water := paint{of: Flat, tint: color.NRGBA{176, 214, 242, 255}}
	b.turn(column{a: vec{}, b: vec{0, 1, 0}, ra: 0.02, rb: 0.018, sides: 8, paint: water})
}}
