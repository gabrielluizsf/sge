package building

import "image/color"

// Morsel is a bite of food, a small round red apple with a stem, held in
// the hand to be eaten. It is not an item the player carries: only the
// person eating shows it.
var Morsel = Building{Name: "morsel", Width: 0.08, Depth: 0.08, build: func(b *builder) {
	skin := paint{of: Flat, tint: color.NRGBA{190, 52, 44, 255}}
	stem := paint{of: Timber}
	b.turn(column{a: vec{0, -0.035, 0}, b: vec{0, 0.035, 0}, ra: 0.036, rb: 0.036, sides: 12, paint: skin})
	b.turn(column{a: vec{0, 0.035, 0}, b: vec{0, 0.055, 0}, ra: 0.006, rb: 0.005, sides: 6, paint: stem})
}}
