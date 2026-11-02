package flora

import (
	"image/color"
	"math"
)

// Cassava is cassava: reddish woody stems topped with hands of long
// narrow leaves, and thick roots at the soil when it is ripe.
var Cassava = Crop{Name: "cassava", Stages: 5,
	leaf: Leaf{Width: 0.28, Widest: 0.5, Veins: 1, Base: color.NRGBA{56, 112, 44, 255}, Tip: color.NRGBA{86, 142, 56, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		wood := g.solid(color.NRGBA{150, 92, 72, 255}, 0.7)
		for s := range 2 + stage/2 {
			yaw := float64(s) * 2.5
			top := vec{math.Sin(yaw) * 0.04, 0.12 + 0.7*g.grown, math.Cos(yaw) * 0.04}
			g.add(stem(rod{vec{}, top, 0.009, wood}))
			g.rosette(whorl{5, top, 0.07 + 0.12*g.grown, 0.25, 0.2, 0.8})
		}
		if stage >= 3 {
			g.hang(vec{0.03, 0.02, 0}, 2, berry{color.NRGBA{190, 150, 100, 255}, 0.035, 0.5})
		}
	}}
