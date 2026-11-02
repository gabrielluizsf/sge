package building

// Stove is a wood-fired kitchen stove, a fogão a lenha: a body of brick
// under an iron cooktop with two rings, the firebox's door and the ash
// door under it, the oven's beside them, a towel rail along its front,
// and its flue rising from the back.
var Stove = Building{Name: "stove", Width: 0.9, Depth: 0.6, build: func(b *builder) { b.stove() }}

// The stove's half width and depth, and its cooktop's height.
const stoveW, stoveD, stoveH = 0.45, 0.3, 0.78

func (b *builder) stove() {
	iron, soot := paint{of: Iron}, paint{of: Soot}
	w, d, h := stoveW, stoveD, stoveH
	b.block(box{lo: vec{-w, 0, -d}, hi: vec{w, h, d}, paint: paint{of: Brick}})
	b.block(box{lo: vec{-w + 0.02, h, -d + 0.02}, hi: vec{w - 0.02, h + 0.03, d - 0.02}, paint: iron})
	for _, x := range []float64{-0.2, 0.12} {
		b.turn(column{a: vec{x, h + 0.03, 0.02}, b: vec{x, h + 0.04, 0.02}, ra: 0.11, rb: 0.1, sides: 16, paint: soot})
	}
	for _, door := range [][4]float64{{-0.36, 0.38, -0.06, 0.62}, {-0.33, 0.12, -0.09, 0.26}, {0.04, 0.3, 0.38, 0.62}} {
		b.block(box{lo: vec{door[0], door[1], d}, hi: vec{door[2], door[3], d + 0.02}, paint: iron})
	}
	for _, x := range []float64{-w + 0.03, w - 0.03} {
		b.block(box{lo: vec{x - 0.01, h - 0.08, d}, hi: vec{x + 0.01, h - 0.04, d + 0.07}, paint: iron})
	}
	b.turn(column{a: vec{-w + 0.03, h - 0.06, d + 0.06}, b: vec{w - 0.03, h - 0.06, d + 0.06}, ra: 0.012, rb: 0.012, sides: 8, paint: iron})
	b.turn(column{a: vec{0.3, h + 0.03, -0.18}, b: vec{0.3, 2.4, -0.18}, ra: 0.065, rb: 0.065, sides: 12, paint: iron})
}
