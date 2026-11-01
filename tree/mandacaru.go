package tree

import "math"

// Mandacaru is the cactus of the semiárido: a ribbed column, and arms
// leaving it low and turning to rise beside it, like a candelabrum, the
// spines pale in its skin. Sapling 0.6 to 1.2 m, young 1.8 to 3, mature
// 3 to 6.
var Mandacaru = Species{Name: "mandacaru", Metres: [Sizes][2]float64{{0.6, 1.2}, {1.8, 3}, {3, 6}}, grow: growMandacaru}

// cactus is a mandacaru being grown: its ribs and skin.
type cactus struct {
	g    *grower
	ribs int
	skin int
	r    float64
}

// growMandacaru grows a mandacaru.
func growMandacaru(g *grower) {
	h := g.size
	c := &cactus{g: g, ribs: g.count(6, 8), skin: g.textured(skin(uint32(g.between(1, 9999))), 0.8), r: 0.035*h + 0.03}
	trunk := g.course(course{vec{}, up, h, 8, 0.01, 0})
	c.limb(trunk, 1)
	arms := [Sizes][2]int{{0, 0}, {1, 3}, {3, 7}}[g.stage]
	turn := g.between(0, 2*math.Pi)
	n := g.count(arms[0], arms[1])
	for i := range n {
		yaw := turn + 2*math.Pi*float64(i)/float64(n) + g.between(-0.4, 0.4)
		out := vec{math.Cos(yaw), 0, math.Sin(yaw)}
		at, _ := along(trunk, g.between(0.25, 0.7))
		start := add(at, scale(out, c.r*0.6))
		path := g.course(course{start, add(out, vec{0, 0.1, 0}), g.between(0.4, 0.62) * h, 12, 0.015, 0.2})
		c.limb(path, 0.62)
	}
}

// limb grows a ribbed column along path, k of the main one's girth.
func (c *cactus) limb(path []vec, k float64) {
	c.g.tube(limb{path: path, r0: c.r * k, r1: c.r * k * 0.8, ribs: c.ribs, rib: 0.22, tiles: c.ribs, span: 0.28, capped: true, mat: c.skin})
}
