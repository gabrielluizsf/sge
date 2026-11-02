package flora

import (
	"image/color"
	"math"
)

// berry is a kind of fruit to hang on a plant: its colour, radius and
// how squashed it is.
type berry struct {
	colour         color.NRGBA
	radius, squash float64
}

// hang grows n of fruit b about the point at, hanging a little below it.
func (g *grower) hang(at vec, n int, b berry) {
	mat := g.solid(b.colour, 0.35)
	for i := range n {
		yaw := float64(i)*2.1 + g.between(0, 1)
		d := vec{math.Sin(yaw) * g.between(0.02, 0.05), -g.between(0.01, 0.05), math.Cos(yaw) * g.between(0.02, 0.05)}
		g.add(fruit(Fruit{Center: add(at, d), Radius: b.radius * g.between(0.85, 1.15), Squash: b.squash}, mat))
	}
}

// climb grows a stem up to a top height tall, leant a little, with n pairs
// of leaves along it of length size, and is its top.
func (g *grower) climb(tall float64, n int, size float64) vec {
	top := vec{g.between(-0.03, 0.03), tall, g.between(-0.03, 0.03)}
	g.add(stem(rod{vec{}, top, 0.008, g.stemMat}))
	for i := range n {
		at := scale(top, (float64(i)+1)/float64(n+1))
		yaw := float64(i)*2.4 + g.between(-0.3, 0.3)
		g.add(leaf(LeafShape{Base: at, Yaw: yaw, Pitch: g.between(0.1, 0.5), Length: size, Width: size * 0.6, Droop: size * 0.3}, 0))
		g.add(leaf(LeafShape{Base: at, Yaw: yaw + math.Pi, Pitch: g.between(0.1, 0.5), Length: size * 0.8, Width: size * 0.5, Droop: size * 0.3}, 0))
	}
	return top
}
