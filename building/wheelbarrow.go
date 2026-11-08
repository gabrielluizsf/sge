package building

// Wheelbarrow is a wooden wheelbarrow heaped with hay: its wheel, shod
// in iron, ahead, along +Z; its handles behind; its tray between, on
// two legs.
var Wheelbarrow = Building{Name: "wheelbarrow", Width: 0.8, Depth: 1.45, build: func(b *builder) {
	wood, iron := paint{of: Timber}, paint{of: Iron}
	b.turn(column{a: vec{-0.04, 0.2, 0.5}, b: vec{0.04, 0.2, 0.5}, ra: 0.2, rb: 0.2, sides: 16, paint: iron})
	b.turn(column{a: vec{-0.05, 0.2, 0.5}, b: vec{0.05, 0.2, 0.5}, ra: 0.17, rb: 0.17, sides: 16, paint: paint{of: Planks}})
	for _, x := range []float64{-0.24, 0.24} {
		b.lay(beam{a: vec{x * 0.4, 0.2, 0.55}, b: vec{x, 0.6, -0.75}, w: 0.05, h: 0.06, paint: wood})
		b.lay(beam{a: vec{x, 0.42, -0.3}, b: vec{x, 0.025, -0.34}, w: 0.05, h: 0.05, paint: wood})
	}
	b.raiseTaper(taper{at: vec{0, 0.3, 0}, h: 0.32, base: [2]float64{0.2, 0.3}, top: [2]float64{0.32, 0.45}, paint: paint{of: Planks}})
	b.roll(ball{c: vec{0, 0.62, 0}, r: 0.4, squash: 0.35, paint: paint{of: Thatch}})
}}
