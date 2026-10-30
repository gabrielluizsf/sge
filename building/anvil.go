package building

// Anvil is a smith's anvil on its stump: about 0.8 m to its face, as a
// smith works at, its horn pointing along X.
var Anvil = Building{Name: "anvil", Width: 0.8, Depth: 0.5, build: func(b *builder) { b.anvil(vec{}) }}

// anvil is an anvil standing at at: a log stump, and on it the iron's
// foot, its waist narrowing and flaring to its face, its heel squared,
// its horn tapering to a point.
func (b *builder) anvil(at vec) {
	iron, wood := paint{of: Iron}, paint{of: Timber}
	b.turn(column{a: at, b: add(at, vec{0, 0.5, 0}), ra: 0.27, rb: 0.25, sides: 10, paint: wood})
	o := func(x, y, z float64) vec { return add(at, vec{x, y, z}) }
	b.block(box{lo: o(-0.17, 0.5, -0.13), hi: o(0.17, 0.57, 0.13), paint: iron})
	b.raiseTaper(taper{at: o(0, 0.57, 0), h: 0.09, base: [2]float64{0.14, 0.11}, top: [2]float64{0.08, 0.055}, paint: iron})
	b.raiseTaper(taper{at: o(0, 0.66, 0), h: 0.07, base: [2]float64{0.08, 0.055}, top: [2]float64{0.21, 0.075}, paint: iron})
	b.block(box{lo: o(-0.22, 0.73, -0.075), hi: o(0.2, 0.83, 0.075), paint: iron})
	b.block(box{lo: o(-0.33, 0.76, -0.055), hi: o(-0.22, 0.83, 0.055), paint: iron})
	b.turn(column{a: o(0.2, 0.785, 0), b: o(0.5, 0.8, 0), ra: 0.06, rb: 0.008, sides: 10, paint: iron})
}
