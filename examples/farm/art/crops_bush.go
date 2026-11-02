package art

import "math"

func cosOf(a float64) float64 { return math.Cos(a) }
func sinOf(a float64) float64 { return math.Sin(a) }

// fruits are the places on a bush where its fruit hangs.
var fruits = [][2]float64{{-3, -4}, {2.6, -6}, {-1, -8}, {3, -3}}

// beansTile is beans: a pole with a vine of leaves, green pods hanging
// when ripe.
func beansTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		tall := 10 + 3*float64(g.stage-2)
		rod{x + 2, y - 0.4, 0, -tall - 2, 0.7, woodBase}.draw(d)
		rod{x, y - 0.4, g.lean * 0.5, -tall, 0.8, leafShade}.draw(d)
		for i := 1.0; i < tall; i += 3 {
			leaf(d, x+g.lean*0.5*i/tall, y-i, [2]float64{3 - 6*float64(int(i)%2), -1.4})
		}
		if g.stage == 3 {
			for _, p := range [][2]float64{{-2, -5}, {3.4, -8}, {-2.4, -8.5}} {
				blob{x + p[0], y + p[1], 0.8, 2.4, "#9ad060", "#5a9a3a"}.draw(d)
			}
		}
	}
	return d.String()
}

// pepperTile is pepper: a small bush with red, tapering peppers.
func pepperTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		sprout(d, 6, 5+float64(g.stage), g.lean)
		if g.stage == 3 {
			for _, p := range fruits[:3] {
				blob{x + p[0], y + p[1], 1, 2.2, "#ff5a40", "#b02a20"}.draw(d)
			}
		} else {
			blob{x + 2, y - 5, 1.1, 1, "#ffffff", "#e8e8d8"}.draw(d)
		}
	}
	return d.String()
}

// coffeeTile is a coffee shrub: dark glossy leaves, white flowers, then
// red cherries.
func coffeeTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		rod{x, y - 0.4, g.lean * 0.4, -7 - 3*float64(g.stage-2), 1.2, woodShade}.draw(d)
		for i, p := range fruits {
			blob{x + p[0]*1.2 + g.lean*0.3, y + p[1] - 1, 3, 1.8, "#2f6b3a", "#1d4a28"}.draw(d)
			switch g.stage {
			case 2:
				blob{x + p[0]*1.2, y + p[1] - 2.4, 0.9, 0.8, "#ffffff", "#e8e8d8"}.draw(d)
			case 3:
				if i != 1 {
					blob{x + p[0]*1.2 + 1, y + p[1] + 0.6, 1, 1, red, redShade}.draw(d)
				}
			}
		}
	}
	return d.String()
}
