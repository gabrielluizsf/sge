package flora

import (
	"image/color"
	"math"
)

// Sugarcane is sugar cane: tall jointed canes, thick, pale green and
// then yellowing, each topped with a plume of long arching leaves.
var Sugarcane = Crop{Name: "sugarcane", Stages: 5,
	leaf: Leaf{Width: 0.3, Widest: 0.3, Veins: 0, Base: color.NRGBA{80, 134, 48, 255}, Tip: color.NRGBA{120, 168, 64, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		colour := color.NRGBA{140, 180, 84, 255}
		if stage == 4 {
			colour = color.NRGBA{196, 184, 90, 255}
		}
		cane := g.solid(colour, 0.5)
		for i := range 3 + stage/2 {
			yaw := float64(i) * 1.9
			base := vec{math.Sin(yaw) * 0.04, 0, math.Cos(yaw) * 0.04}
			top := add(base, vec{g.between(-0.02, 0.02), 0.2 + 1.5*g.grown, g.between(-0.02, 0.02)})
			g.add(stem(rod{base, top, 0.014, cane}))
			g.rosette(whorl{4, top, 0.15 + 0.3*g.grown, 0.1, 0.5, 1.0})
		}
	}}
