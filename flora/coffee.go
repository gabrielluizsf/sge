package flora

import "image/color"

// Coffee is a coffee shrub: a woody trunk with glossy, dark green leaves
// in pairs, white flowers, then clusters of cherries that go red.
var Coffee = Crop{Name: "coffee", Stages: 5,
	leaf: Leaf{Width: 0.5, Widest: 0.45, Veins: 6, Base: color.NRGBA{30, 84, 40, 255}, Tip: color.NRGBA{52, 112, 52, 255}, Vein: color.NRGBA{110, 164, 100, 255}},
	grow: func(g *grower, stage int) {
		top := g.climb(0.12+0.8*g.grown, 4+stage*3, 0.07+0.08*g.grown)
		switch {
		case stage == 3:
			g.hang(scale(top, 0.6), 5, berry{color.NRGBA{250, 250, 244, 255}, 0.01, 0.7})
		case stage == 4:
			g.hang(scale(top, 0.55), 8, berry{color.NRGBA{190, 24, 40, 255}, 0.014, 1})
		}
	}}
