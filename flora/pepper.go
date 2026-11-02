package flora

import "image/color"

// Pepper is a pepper plant: a small bush of smooth, pointed leaves with
// peppers hanging from it, green, then red.
var Pepper = Crop{Name: "pepper", Stages: 5,
	leaf: Leaf{Width: 0.5, Widest: 0.4, Veins: 4, Base: color.NRGBA{58, 116, 44, 255}, Tip: color.NRGBA{88, 148, 58, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		top := g.climb(0.08+0.4*g.grown, 3+stage*2, 0.06+0.06*g.grown)
		if stage < 3 {
			return
		}
		c := color.NRGBA{96, 150, 56, 255}
		if stage == 4 {
			c = color.NRGBA{218, 40, 28, 255}
		}
		g.hang(scale(top, 0.6), 2+stage-3, berry{c, 0.016, 2.4})
	}}
