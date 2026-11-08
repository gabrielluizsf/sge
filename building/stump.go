package building

import "math"

// Stump is what is left of a felled tree: a short round of trunk in its
// bark, sawn level, its roots flaring into the ground.
var Stump = Building{Name: "stump", Width: 0.69, Depth: 0.69, build: func(b *builder) { b.stump() }}

func (b *builder) stump() {
	bark := paint{of: Bark}
	b.turn(column{a: vec{0, 0, 0}, b: vec{0, 0.36, 0}, ra: 0.27, rb: 0.23, sides: 16, paint: bark})
	b.turn(column{a: vec{0, 0.36, 0}, b: vec{0, 0.38, 0}, ra: 0.225, rb: 0.225, sides: 16, paint: paint{of: Timber}})
	for i := range 5 {
		a := 2*math.Pi*float64(i)/5 + 0.4
		c, s := math.Cos(a), math.Sin(a)
		b.turn(column{a: vec{0.16 * c, 0.14, 0.16 * s}, b: vec{0.36 * c, 0.03, 0.36 * s}, ra: 0.07, rb: 0.025, sides: 8, paint: bark})
	}
}
