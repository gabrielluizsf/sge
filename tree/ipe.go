package tree

import (
	"image/color"

	"github.com/gabrielluizsf/sge/flora"
)

// Ipe is the ipê-amarelo of the campo: a crooked trunk forking into
// limbs that carry a round crown, bare of leaves and covered in yellow
// flowers. Sapling 1.5 to 2.5 m, young 4 to 6, mature 8 to 12.
var Ipe = Species{Name: "ipe", Metres: [Sizes][2]float64{{1.5, 2.5}, {4, 6}, {8, 12}}, grow: growIpe(true)}

// IpeGreen is the ipê out of flower: the same crooked tree in green
// leaves.
var IpeGreen = Species{Name: "ipe_green", Metres: Ipe.Metres, grow: growIpe(false)}

// growIpe grows an ipê, in flower or in leaf.
func growIpe(flowering bool) func(g *grower) {
	return func(g *grower) {
		h := g.size
		cover := foliage{n: [2]int{10, 15}, size: [2]float64{0.08 * h, 0.13 * h}, aspect: 1, reach: 0.5, centre: vec{0, 0.7 * h, 0}}
		if flowering {
			cover.mat = g.cutout("blossom", blossoms(uint32(g.between(1, 9999)), color.NRGBA{250, 206, 36, 255}, color.NRGBA{196, 112, 20, 255}))
		} else {
			cover.mat = g.cutout("leaf", flora.Leaf{Width: 0.4, Widest: 0.5, Teeth: 0.3, Veins: 6, Base: color.NRGBA{70, 120, 48, 255},
				Tip: color.NRGBA{98, 150, 62, 255}, Vein: color.NRGBA{170, 204, 120, 255}}.Texture(96, uint32(g.between(1, 9999))))
			cover.size = [2]float64{0.045 * h, 0.075 * h}
			cover.n = [2]int{18, 26}
			cover.aspect = 0.7
		}
		c := &crown{depth: 4, kids: [2]int{2, 3}, fork: [2]int{3, 5}, spread: [2]float64{0.4, 0.85}, shrink: 0.68, thin: 0.6, wobble: 0.16, lift: 0.04,
			flare: 0.5, sides: 7, bark: g.bark(g.barks.Smooth), span: g.barks.Smooth.Span, foliage: cover}
		if g.stage == Sapling {
			c.depth = 3
		}
		lean := vec{g.between(-0.35, 0.35), 1, g.between(-0.35, 0.35)}
		g.grow(c, bough{from: vec{}, dir: lean, length: 0.5 * h, r: 0.022*h + 0.01})
	}
}
