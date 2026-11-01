package tree

import (
	"image/color"

	"github.com/gabrielluizsf/sge/flora"
)

// Cashew is the cajueiro of the campo and the coast: a short, crooked
// trunk and a low crown sprawling wider than it is tall, its branches
// reaching out and drooping, in broad leaves. Measured across: sapling
// 1 to 2 m, young 3 to 5, mature 6 to 8.
var Cashew = Species{Name: "cashew", Metres: [Sizes][2]float64{{1, 2}, {3, 5}, {6, 8}}, Wide: true, grow: growCashew}

// growCashew grows a cajueiro.
func growCashew(g *grower) {
	w := g.size
	leaf := flora.Leaf{Width: 0.46, Widest: 0.72, Veins: 9, Base: color.NRGBA{78, 124, 50, 255}, Tip: color.NRGBA{58, 106, 46, 255}, Vein: color.NRGBA{156, 190, 112, 255}}
	c := &crown{depth: 4, kids: [2]int{2, 4}, fork: [2]int{3, 4}, spread: [2]float64{0.35, 0.8}, shrink: 0.62, thin: 0.6, wobble: 0.14, lift: -0.02,
		flare: 0.5, sides: 7, bark: g.bark(g.barks.Smooth), span: g.barks.Smooth.Span, foliage: foliage{
			mat: g.cutout("leaf", leaf.Texture(96, uint32(g.between(1, 9999)))), n: [2]int{10, 15},
			size: [2]float64{0.05 * w, 0.08 * w}, aspect: 0.62, reach: 0.7, centre: vec{0, -0.1 * w, 0}}}
	if g.stage == Sapling {
		c.depth = 3
	}
	lean := vec{g.between(-0.5, 0.5), 1, g.between(-0.5, 0.5)}
	g.grow(c, bough{from: vec{}, dir: lean, length: 0.22 * w, r: 0.022 * w, depth: 0})
}
