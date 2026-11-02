package flora

import "image/color"

// Strawberry is a strawberry: a low mat of toothed, three-lobed leaves,
// with white flowers, then red berries on top of the plant.
var Strawberry = Crop{Name: "strawberry", Stages: 5,
	leaf: Leaf{Width: 0.62, Widest: 0.55, Teeth: 0.8, Veins: 5, Base: color.NRGBA{58, 116, 42, 255}, Tip: color.NRGBA{90, 150, 56, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		g.rosette(whorl{4 + stage*2, vec{}, 0.05 + 0.1*g.grown, 0.9, 0.5, 1.2})
		switch {
		case stage == 3:
			g.hang(vec{0, 0.05, 0}, 3, berry{color.NRGBA{250, 250, 245, 255}, 0.01, 0.6})
		case stage == 4:
			g.hang(vec{0, 0.06, 0}, 5, berry{color.NRGBA{214, 36, 44, 255}, 0.016, 1.1})
		}
	}}
