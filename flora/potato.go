package flora

import "image/color"

// Potato is a potato: a bushy plant of divided leaves, flowering pale as
// it ripens, a tuber showing at its foot.
var Potato = Crop{Name: "potato", Stages: 5,
	leaf: Leaf{Width: 0.5, Widest: 0.5, Teeth: 0.5, Veins: 5, Base: color.NRGBA{62, 112, 44, 255}, Tip: color.NRGBA{92, 144, 58, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		g.rosette(whorl{4 + stage*3, vec{0, 0.02, 0}, 0.06 + 0.2*g.grown, 0.7, 0.6, 1.3})
		if stage < 3 {
			return
		}
		g.hang(vec{0, 0.22, 0}, 3+stage-3, berry{color.NRGBA{240, 236, 250, 255}, 0.012, 0.8})
		g.hang(vec{0.04, 0.01, 0.02}, 2, berry{color.NRGBA{200, 160, 100, 255}, 0.03, 0.75})
	}}
