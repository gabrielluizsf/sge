package flora

import "image/color"

// Beans are climbing beans: a vine up a pole, heart-shaped leaves, and long
// green pods hanging once it ripens.
var Beans = Crop{Name: "beans", Stages: 5,
	leaf: Leaf{Width: 0.7, Widest: 0.4, Veins: 3, Base: color.NRGBA{66, 124, 46, 255}, Tip: color.NRGBA{96, 156, 60, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		h := 0.1 + 0.8*g.grown
		pole := g.solid(color.NRGBA{150, 112, 66, 255}, 0.8)
		g.add(stem(rod{vec{0.03, 0, 0}, vec{0.03, h + 0.1, 0}, 0.007, pole}))
		top := g.climb(h, 2+stage*2, 0.06+0.06*g.grown)
		if stage < 3 {
			return
		}
		pods := g.solid(color.NRGBA{120, 176, 70, 255}, 0.5)
		for range 2 + stage - 3 {
			at := scale(top, g.between(0.35, 0.8))
			g.add(stem(rod{at, add(at, vec{g.between(-0.02, 0.02), -0.1, 0.02}), 0.006, pods}))
		}
	}}
