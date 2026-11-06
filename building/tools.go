package building

import "image/color"

// Items are the things a player carries, as small models to bake into
// icons: laid on the ground, long along X, to be seen from a corner. The
// hoe, the watering can, the logs, a coin, an axe with each head and
// the pickaxe.
var Items = []Building{Hoe, Can, Logs, Coin, AxeStone, AxeIron, AxeGold, AxeDiamond, Pickaxe}

// Hoe is a long wooden handle with an iron blade across its end.
var Hoe = Building{Name: "hoe", Width: 1, Depth: 0.3, build: func(b *builder) {
	wood, iron := paint{of: Timber}, paint{of: Iron}
	b.lay(beam{a: vec{-0.55, 0.03, 0}, b: vec{0.42, 0.03, 0}, w: 0.04, h: 0.04, paint: wood})
	b.turn(column{a: vec{0.28, 0.03, 0}, b: vec{0.42, 0.03, 0}, ra: 0.035, rb: 0.035, sides: 10, paint: iron})
	b.block(box{lo: vec{0.4, 0.005, -0.14}, hi: vec{0.48, 0.05, 0.14}, paint: iron})
}}

// Can is a watering can: a round body, a long spout with its rose, and a
// handle behind.
var Can = Building{Name: "can", Width: 0.7, Depth: 0.3, build: func(b *builder) {
	tin := paint{of: Flat, tint: color.NRGBA{70, 118, 168, 255}}
	rim := paint{of: Iron}
	b.turn(column{a: vec{0, 0, 0}, b: vec{0, 0.3, 0}, ra: 0.15, rb: 0.15, sides: 20, paint: tin})
	b.turn(column{a: vec{0, 0.3, 0}, b: vec{0, 0.325, 0}, ra: 0.155, rb: 0.155, sides: 20, paint: rim})
	b.lay(beam{a: vec{0.1, 0.07, 0}, b: vec{0.46, 0.34, 0}, w: 0.05, h: 0.05, paint: tin})
	b.turn(column{a: vec{0.44, 0.33, 0}, b: vec{0.55, 0.39, 0}, ra: 0.025, rb: 0.075, sides: 14, paint: rim})
	b.lay(beam{a: vec{-0.12, 0.3, 0}, b: vec{-0.26, 0.2, 0}, w: 0.04, h: 0.05, paint: rim})
	b.lay(beam{a: vec{-0.26, 0.2, 0}, b: vec{-0.13, 0.05, 0}, w: 0.04, h: 0.05, paint: rim})
}}

// Logs are three logs of firewood, two with a third over them.
var Logs = Building{Name: "wood", Width: 0.6, Depth: 0.4, build: func(b *builder) {
	wood := paint{of: Timber}
	for _, at := range []vec{{0, 0.075, -0.085}, {0.03, 0.075, 0.085}, {0.015, 0.2, 0}} {
		b.turn(column{a: vec{at[0] - 0.27, at[1], at[2]}, b: vec{at[0] + 0.27, at[1], at[2]}, ra: 0.075, rb: 0.075, sides: 14, paint: wood})
	}
}}

// Coin is a gold coin, a thin round disc with a raised rim.
var Coin = Building{Name: "coin", Width: 0.36, Depth: 0.36, build: func(b *builder) {
	gold := paint{of: Flat, tint: color.NRGBA{236, 186, 58, 255}}
	rim := paint{of: Flat, tint: color.NRGBA{250, 214, 100, 255}}
	b.turn(column{a: vec{0, 0, 0}, b: vec{0, 0.03, 0}, ra: 0.17, rb: 0.17, sides: 28, paint: rim})
	b.turn(column{a: vec{0, 0.03, 0}, b: vec{0, 0.036, 0}, ra: 0.13, rb: 0.13, sides: 28, paint: gold})
}}

// The axes, by what their heads are made of.
var (
	AxeStone   = axeOf("axe_stone", paint{of: Stone})
	AxeIron    = axeOf("axe_iron", paint{of: Iron})
	AxeGold    = axeOf("axe_gold", paint{of: Flat, tint: color.NRGBA{234, 184, 54, 255}})
	AxeDiamond = axeOf("axe_diamond", paint{of: Flat, tint: color.NRGBA{120, 222, 236, 255}})
)

// axeOf is an axe called name, its head of head: a wooden handle, an eye
// round it, and a blade that flares to an edge.
func axeOf(name string, head paint) Building {
	return Building{Name: name, Width: 0.9, Depth: 0.4, build: func(b *builder) {
		b.lay(beam{a: vec{-0.5, 0.03, 0}, b: vec{0.36, 0.03, 0}, w: 0.045, h: 0.045, paint: paint{of: Timber}})
		b.block(box{lo: vec{0.2, 0.005, -0.045}, hi: vec{0.33, 0.07, 0.045}, paint: head})
		b.convex(hexa([8]vec{
			{0.2, 0.015, 0.045}, {0.33, 0.015, 0.045}, {0.4, 0.015, 0.21}, {0.13, 0.015, 0.21},
			{0.2, 0.06, 0.045}, {0.33, 0.06, 0.045}, {0.4, 0.06, 0.21}, {0.13, 0.06, 0.21},
		}, head))
	}}
}

// Pickaxe is a pickaxe: a wooden handle and, across its end, an iron head
// of two points curving away from an eye, one to each side. Held two-handed
// it is swung as the axe is, the points facing the rock.
var Pickaxe = Building{Name: "pickaxe", Width: 0.9, Depth: 0.7, build: func(b *builder) {
	iron := paint{of: Iron}
	b.lay(beam{a: vec{-0.5, 0.03, 0}, b: vec{0.36, 0.03, 0}, w: 0.045, h: 0.045, paint: paint{of: Timber}})
	b.block(box{lo: vec{0.2, 0.005, -0.05}, hi: vec{0.33, 0.07, 0.05}, paint: iron})
	for _, side := range []float64{1, -1} {
		b.convex(hexa([8]vec{
			{0.2, 0.015, 0.05 * side}, {0.33, 0.015, 0.05 * side}, {0.3, 0.015, 0.3 * side}, {0.27, 0.015, 0.3 * side},
			{0.2, 0.06, 0.05 * side}, {0.33, 0.06, 0.05 * side}, {0.3, 0.06, 0.3 * side}, {0.27, 0.06, 0.3 * side},
		}, iron))
	}
}}
