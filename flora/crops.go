package flora

import (
	"image/color"
	"math"
)

var (
	green     = color.NRGBA{70, 128, 48, 255}
	lightLeaf = color.NRGBA{150, 196, 92, 255}
)

// Carrot is a carrot: a tuft of fine, feathery leaves, and at the last
// stages its orange shoulder showing at the soil.
var Carrot = Crop{Name: "carrot", Stages: 5,
	leaf: Leaf{Width: 0.46, Widest: 0.55, Teeth: 0.9, Veins: 7, Base: color.NRGBA{60, 112, 40, 255}, Tip: color.NRGBA{98, 150, 58, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		g.rosette(whorl{3 + stage*2, vec{}, 0.04 + 0.24*g.grown, 0.5, 0.9, 1.35})
		if stage >= 3 {
			orange := g.solid(color.NRGBA{226, 118, 34, 255}, 0.5)
			r := 0.012 + 0.006*float64(stage-2)
			g.add(fruit(Fruit{Center: vec{0, r * 0.4, 0}, Radius: r, Squash: 1}, orange))
		}
	}}

// Lettuce is a lettuce: broad, frilled, pale leaves in a head that
// closes as it grows.
var Lettuce = Crop{Name: "lettuce", Stages: 5,
	leaf: Leaf{Width: 0.46, Widest: 0.62, Teeth: 0.25, Veins: 6, Base: color.NRGBA{170, 206, 110, 255}, Tip: color.NRGBA{108, 166, 66, 255}, Vein: color.NRGBA{214, 236, 170, 255}},
	grow: func(g *grower, stage int) {
		g.rosette(whorl{5 + stage*3, vec{}, 0.05 + 0.14*g.grown, 0.95, 0.25, 1.25})
	}}

// Tomato is a tomato: a stem climbing with its leaves, and clusters of
// fruit, green before they ripen red.
var Tomato = Crop{Name: "tomato", Stages: 5,
	leaf: Leaf{Width: 0.46, Widest: 0.45, Teeth: 0.35, Veins: 5, Base: green, Tip: color.NRGBA{84, 140, 58, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		h := 0.1 + 0.75*g.grown
		top := vec{g.between(-0.03, 0.03), h, g.between(-0.03, 0.03)}
		g.add(stem(rod{vec{}, top, 0.008, g.stemMat}))
		nodes := 2 + stage*2
		for i := range nodes {
			at := scale(top, (float64(i)+1)/float64(nodes+1))
			yaw := float64(i)*2.4 + g.between(-0.3, 0.3)
			l := 0.08 + 0.1*g.grown
			g.add(leaf(LeafShape{Base: at, Yaw: yaw, Pitch: g.between(0.1, 0.5), Length: l, Width: l * 0.6, Droop: l * 0.35}, 0))
			g.add(leaf(LeafShape{Base: at, Yaw: yaw + math.Pi + g.between(-0.4, 0.4), Pitch: g.between(0.1, 0.5), Length: l * 0.8, Width: l * 0.5, Droop: l * 0.3}, 0))
		}
		if stage < 3 {
			return
		}
		red := color.NRGBA{205, 40, 30, 255}
		if stage == 3 {
			red = color.NRGBA{120, 165, 60, 255}
		}
		fruitMat := g.solid(red, 0.3)
		for c := range 2 + stage - 3 {
			at := scale(top, g.between(0.35, 0.75))
			yaw := float64(c)*2.1 + g.between(0, 1)
			for range 3 {
				d := vec{math.Sin(yaw) * g.between(0.03, 0.06), -g.between(0.02, 0.06), math.Cos(yaw) * g.between(0.03, 0.06)}
				g.add(fruit(Fruit{Center: add(at, d), Radius: g.between(0.018, 0.028), Squash: 0.85}, fruitMat))
			}
		}
	}}

// Corn is maize: a tall stalk of long arching leaves, a tassel at its
// top and ears along it as it ripens.
var Corn = Crop{Name: "corn", Stages: 5,
	leaf: Leaf{Width: 0.44, Widest: 0.35, Veins: 0, Base: color.NRGBA{66, 120, 44, 255}, Tip: color.NRGBA{104, 150, 60, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		h := 0.15 + 1.75*g.grown
		top := vec{0, h, 0}
		g.add(stem(rod{vec{}, top, 0.012 + 0.006*g.grown, g.stemMat}))
		n := 3 + stage*2
		for i := range n {
			at := scale(top, (float64(i)+0.5)/float64(n+1))
			l := (0.25 + 0.55*g.grown) * g.between(0.85, 1.1)
			g.add(leaf(LeafShape{Base: at, Yaw: float64(i)*math.Pi + g.between(-0.4, 0.4), Pitch: g.between(0.5, 0.9), Length: l, Width: 0.085, Droop: l * 0.45}, 0))
		}
		if stage < 3 {
			return
		}
		tassel := g.solid(color.NRGBA{196, 170, 92, 255}, 0.8)
		for i := range 5 {
			yaw := float64(i) * 1.26
			g.add(stem(rod{top, add(top, vec{math.Sin(yaw) * 0.05, 0.12, math.Cos(yaw) * 0.05}), 0.002, tassel}))
		}
		husk := g.solid(color.NRGBA{150, 170, 80, 255}, 0.7)
		if stage == 4 {
			husk = g.solid(color.NRGBA{190, 176, 98, 255}, 0.75)
		}
		at := scale(top, 0.55)
		g.add(stem(rod{at, add(at, vec{0.06, 0.16, 0.02}), 0.022, husk}))
	}}

// Wheat is a clump of wheat: thin stalks, green, then golden, each with
// its head of grain.
var Wheat = Crop{Name: "wheat", Stages: 5,
	leaf: Leaf{Width: 0.42, Widest: 0.3, Base: color.NRGBA{92, 138, 56, 255}, Tip: color.NRGBA{120, 160, 70, 255}, Vein: lightLeaf},
	grow: func(g *grower, stage int) {
		stalkColour := color.NRGBA{110, 150, 62, 255}
		headColour := color.NRGBA{130, 165, 70, 255}
		if stage == 4 {
			stalkColour, headColour = color.NRGBA{214, 184, 106, 255}, color.NRGBA{226, 190, 110, 255}
		}
		stalk := g.solid(stalkColour, 0.8)
		head := g.solid(headColour, 0.75)
		for range 8 + stage*4 {
			base := vec{g.between(-0.06, 0.06), 0, g.between(-0.06, 0.06)}
			h := (0.08 + 0.75*g.grown) * g.between(0.85, 1.1)
			top := add(base, vec{g.between(-0.08, 0.08) * g.grown, h, g.between(-0.08, 0.08) * g.grown})
			g.add(stem(rod{base, top, 0.002, stalk}))
			if stage >= 3 {
				g.add(fruit(Fruit{Center: add(top, vec{0, 0.03, 0}), Radius: 0.008, Squash: 4}, head))
			} else {
				l := h * 0.6
				g.add(leaf(LeafShape{Base: base, Yaw: g.between(0, 6.3), Pitch: 1.1, Length: l, Width: 0.014, Droop: l * 0.2}, 0))
			}
		}
	}}

// Pumpkin is a pumpkin: broad leaves along vines on the ground, and its
// ribbed fruit swelling, green, then orange.
var Pumpkin = Crop{Name: "pumpkin", Stages: 5,
	leaf: Leaf{Width: 0.5, Widest: 0.5, Teeth: 0.4, Veins: 4, Base: color.NRGBA{58, 108, 40, 255}, Tip: color.NRGBA{80, 130, 50, 255}, Vein: color.NRGBA{160, 196, 120, 255}},
	grow: func(g *grower, stage int) {
		vines := 2 + stage/2
		for v := range vines {
			yaw := float64(v)*2.4 + g.between(-0.3, 0.3)
			reach := 0.1 + 0.4*g.grown
			end := vec{math.Sin(yaw) * reach, 0.01, math.Cos(yaw) * reach}
			g.add(stem(rod{vec{0, 0.01, 0}, end, 0.006, g.stemMat}))
			for k := range 2 + stage {
				at := scale(end, (float64(k)+1)/float64(2+stage))
				l := 0.08 + 0.16*g.grown
				g.add(leaf(LeafShape{Base: at, Yaw: yaw + g.between(-1.2, 1.2), Pitch: g.between(0.4, 0.8), Length: l, Width: l * 0.95, Droop: l * 0.25}, 0))
			}
		}
		if stage < 3 {
			return
		}
		c := color.NRGBA{226, 120, 30, 255}
		if stage == 3 {
			c = color.NRGBA{140, 150, 60, 255}
		}
		r := 0.07 + 0.05*float64(stage-3)
		g.add(fruit(Fruit{Center: vec{0.05, r * 0.7, 0.02}, Radius: r, Squash: 0.75, Rib: 0.12, Ribs: 10}, g.solid(c, 0.45)))
	}}

// Crops are every crop flora grows.
var Crops = []Crop{Carrot, Lettuce, Tomato, Corn, Wheat, Pumpkin}
