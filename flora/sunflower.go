package flora

import (
	"image/color"
	"math"
)

// Sunflower is a sunflower: a tall rough stem with big heart-shaped
// leaves, a green bud, then a broad head of yellow petals round a dark
// brown disc.
var Sunflower = Crop{Name: "sunflower", Stages: 5,
	leaf: Leaf{Width: 0.8, Widest: 0.35, Teeth: 0.3, Veins: 5, Base: color.NRGBA{64, 118, 44, 255}, Tip: color.NRGBA{94, 150, 58, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		top := g.climb(0.12+1.2*g.grown, 2+stage, 0.1+0.12*g.grown)
		switch {
		case stage == 3:
			g.add(fruit(Fruit{Center: top, Radius: 0.04, Squash: 0.8}, g.solid(color.NRGBA{90, 140, 56, 255}, 0.6)))
		case stage == 4:
			petals := g.solid(color.NRGBA{250, 200, 30, 255}, 0.5)
			for i := range 12 {
				a := float64(i) * math.Pi / 6
				g.add(fruit(Fruit{Center: add(top, vec{math.Sin(a) * 0.1, 0, math.Cos(a) * 0.1}), Radius: 0.04, Squash: 0.3}, petals))
			}
			g.add(fruit(Fruit{Center: top, Radius: 0.085, Squash: 0.3}, g.solid(color.NRGBA{78, 48, 26, 255}, 0.8)))
		}
	}}
