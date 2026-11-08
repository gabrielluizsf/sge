package flora

import "image/color"

// Withered is a plant that went a day dry: a bent brown stem and limp,
// curled, dried leaves lying about it.
var Withered = Crop{Name: "withered", Stages: 1,
	leaf: Leaf{Width: 0.4, Widest: 0.5, Veins: 3, Base: color.NRGBA{140, 112, 60, 255}, Tip: color.NRGBA{170, 140, 76, 255}, Vein: color.NRGBA{190, 164, 100, 255}},
	grow: func(g *grower, stage int) {
		dry := g.solid(color.NRGBA{128, 100, 56, 255}, 0.9)
		mid := vec{0.02, 0.14, 0}
		g.add(stem(rod{vec{}, mid, 0.008, dry}))
		g.add(stem(rod{mid, vec{0.09, 0.1, 0.02}, 0.006, dry}))
		for i := range 5 {
			yaw := float64(i)*1.3 + g.between(-0.4, 0.4)
			g.add(leaf(LeafShape{Base: vec{0, 0.03 + 0.02*float64(i%3), 0}, Yaw: yaw, Pitch: 0.1, Length: g.between(0.08, 0.13), Width: 0.05, Droop: 0.09}, 0))
		}
	}}
