package flora

import (
	"image/color"
	"math"
)

// Watermelon is a watermelon: long vines of deeply cut leaves along the
// ground, and a big, dark green, ribbed melon swelling on it.
var Watermelon = Crop{Name: "watermelon", Stages: 5,
	leaf: Leaf{Width: 0.55, Widest: 0.5, Teeth: 0.9, Veins: 4, Base: color.NRGBA{60, 112, 44, 255}, Tip: color.NRGBA{88, 140, 56, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		for v := range 2 + stage/2 {
			yaw := float64(v)*2.2 + g.between(-0.3, 0.3)
			reach := 0.1 + 0.5*g.grown
			end := vec{math.Sin(yaw) * reach, 0.01, math.Cos(yaw) * reach}
			g.add(stem(rod{vec{0, 0.01, 0}, end, 0.005, g.stemMat}))
			for k := range 2 + stage {
				at := scale(end, (float64(k)+1)/float64(2+stage))
				l := 0.07 + 0.13*g.grown
				g.add(leaf(LeafShape{Base: at, Yaw: yaw + g.between(-1.2, 1.2), Pitch: g.between(0.4, 0.8), Length: l, Width: l * 0.9, Droop: l * 0.25}, 0))
			}
		}
		if stage < 3 {
			return
		}
		r := 0.05 + 0.06*float64(stage-3)
		c := color.NRGBA{52, 128, 60, 255}
		g.add(fruit(Fruit{Center: vec{0.06, r * 0.8, 0.03}, Radius: r, Squash: 0.85, Rib: 0.06, Ribs: 12}, g.solid(c, 0.4)))
	}}
