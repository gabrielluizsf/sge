package tree

import (
	"image/color"
	"math"
)

// Coconut is the coqueiro of the coast: a slender trunk, ringed and
// curving as it leans, and at its top a crown of long pinnate fronds
// arching out and down over a bunch of coconuts. Sapling 2 to 3 m,
// young 5 to 8, mature 10 to 15.
var Coconut = Species{Name: "coconut", Metres: [Sizes][2]float64{{2, 3}, {5, 8}, {10, 15}}, grow: growCoconut}

// growCoconut grows a coqueiro.
func growCoconut(g *grower) {
	h := g.size
	trunkH := h * [Sizes]float64{0.35, 0.78, 0.82}[g.stage]
	r := 0.014*h + 0.02
	lean := vec{g.between(-0.5, 0.5), 1, g.between(-0.5, 0.5)}
	trunk := g.course(course{vec{}, lean, trunkH, 14, 0.012, 0.035})
	g.tube(limb{path: trunk, r0: r, r1: r * 0.72, flare: 0.9, sides: 10, mat: g.bark(g.barks.Palm), span: g.barks.Palm.Span})
	top, _ := along(trunk, 1)
	crown := g.fronds(top, trunkH)
	if g.stage == Sapling {
		return
	}
	nuts := g.solid(color.NRGBA{112, 128, 52, 255}, 0.5)
	for range g.count(4, 9) {
		a := g.between(0, 2*math.Pi)
		d := g.between(0.1, 0.3) * crown
		g.ball(ball{at: add(top, vec{math.Cos(a) * d, -g.between(0.05, 0.3) * crown, math.Sin(a) * d}), r: 0.014*h + 0.04, squash: 1.2, mat: nuts})
	}
}

// fronds grows the crown of fronds at top, each as long as a share of
// the trunk's height says; the length of a frond.
func (g *grower) fronds(top vec, trunkH float64) float64 {
	n := g.count(13, 17)
	if g.stage == Sapling {
		n = g.count(6, 8)
	} else if g.stage == Young {
		n = g.count(10, 13)
	}
	long := math.Max(0.36*trunkH, 0.4*g.size)
	if g.stage == Sapling {
		long = 0.4 * g.size
	}
	mat := g.cutout("frond", frond(uint32(g.between(1, 9999))))
	turn := g.between(0, 2*math.Pi)
	for i := range n {
		rel := float64(i) / float64(n)
		yaw := turn + float64(i)*2.39996
		pitch := 1.15 - 1.4*rel + g.between(-0.1, 0.1)
		dir := vec{math.Cos(yaw) * math.Cos(pitch), math.Sin(pitch), math.Sin(yaw) * math.Cos(pitch)}
		l := long * g.between(0.85, 1.1)
		path := g.course(course{top, dir, l, 9, 0.02, -0.1})
		g.ribbon(ribbon{path: path, width: l * 0.3, fold: 0.55, mat: mat})
	}
	return long
}
