package building

import "image/color"

// MilkStool is what one milks with: a low three-legged stool under them
// and a pail forward of it, under the teats. It is no thing in the world:
// only the person milking shows it.
var MilkStool = Building{Name: "milk_stool", Width: 0.3, Depth: 0.9, build: func(b *builder) {
	wood, metal := paint{of: Timber}, paint{of: Flat, tint: color.NRGBA{150, 152, 150, 255}}
	const seat = 0.28
	b.turn(column{a: vec{0, seat - 0.035, 0}, b: vec{0, seat, 0}, ra: 0.13, rb: 0.13, sides: 12, paint: paint{of: Planks}})
	for _, foot := range []vec{{0.12, 0, 0.08}, {-0.12, 0, 0.08}, {0, 0, -0.13}} {
		b.turn(column{a: foot, b: vec{foot[0] / 2, seat - 0.03, foot[2] / 2}, ra: 0.018, rb: 0.018, sides: 6, paint: wood})
	}
	const pail = 0.5 // forward of the stool
	b.turn(column{a: vec{0, 0, pail}, b: vec{0, 0.27, pail}, ra: 0.1, rb: 0.135, sides: 14, paint: metal})
	b.turn(column{a: vec{0, 0.262, pail}, b: vec{0, 0.268, pail}, ra: 0.122, rb: 0.122, sides: 14, paint: paint{of: Flat, tint: color.NRGBA{244, 240, 228, 255}}})
}}

// MilkBottle is a bottle of milk, an item carried: glass full of white,
// corked.
var MilkBottle = Building{Name: "milk", Width: 0.12, Depth: 0.12, build: func(b *builder) {
	white := paint{of: Flat, tint: color.NRGBA{244, 240, 228, 255}}
	b.turn(column{a: vec{0, 0, 0}, b: vec{0, 0.2, 0}, ra: 0.055, rb: 0.055, sides: 14, paint: white})
	b.turn(column{a: vec{0, 0.2, 0}, b: vec{0, 0.27, 0}, ra: 0.055, rb: 0.024, sides: 14, paint: white})
	b.turn(column{a: vec{0, 0.27, 0}, b: vec{0, 0.31, 0}, ra: 0.024, rb: 0.024, sides: 10, paint: white})
	b.turn(column{a: vec{0, 0.31, 0}, b: vec{0, 0.335, 0}, ra: 0.027, rb: 0.024, sides: 10, paint: paint{of: Timber}})
}}

// HandBottle is the bottle of milk as it is held to drink from: a size for the
// hand, its neck narrowing to the lips, unstopped. Only the person drinking
// shows it.
var HandBottle = Building{Name: "hand_bottle", Width: 0.08, Depth: 0.08, build: func(b *builder) {
	glass := paint{of: Flat, tint: color.NRGBA{244, 240, 228, 255}}
	b.turn(column{a: vec{0, 0, 0}, b: vec{0, 0.13, 0}, ra: 0.036, rb: 0.036, sides: 12, paint: glass})
	b.turn(column{a: vec{0, 0.13, 0}, b: vec{0, 0.175, 0}, ra: 0.036, rb: 0.016, sides: 12, paint: glass})
	b.turn(column{a: vec{0, 0.175, 0}, b: vec{0, 0.2, 0}, ra: 0.016, rb: 0.017, sides: 8, paint: paint{of: Flat, tint: color.NRGBA{210, 205, 190, 255}}})
}}
