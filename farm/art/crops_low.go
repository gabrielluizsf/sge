package art

// lettuceTile is a lettuce: a tuft, then a head of pale frilled leaves.
func lettuceTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage == 2:
		sprout(d, 5, 4, g.lean)
	case g.stage == 3:
		for i, dx := range []float64{-3.4, 3.4, 0, -1.8, 1.8} {
			blob{x + dx + g.lean*0.3, y - 2.4 - float64(i%3)*1.4, 3.2, 2.6, "#c8ec8c", "#7fc65a"}.draw(d)
		}
		blob{x + g.lean*0.5, y - 5.6, 2.6, 2.2, "#e3f7b0", "#a4d672"}.draw(d)
	}
	return d.String()
}

// potatoTile is a potato: a bushy plant, flowering when ripe, a tuber at
// its foot.
func potatoTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		sprout(d, 6, 5+float64(g.stage), g.lean)
		if g.stage == 3 {
			for _, p := range [][2]float64{{-2.4, -6}, {2.2, -7.5}, {0, -9}} {
				blob{x + p[0] + g.lean*0.6, y + p[1], 1.1, 1, "#ffffff", "#e8c8f0"}.draw(d)
			}
			blob{x + 2.8, y - 0.6, 1.9, 1.3, "#e0b878", "#b48a4e"}.draw(d)
		}
	}
	return d.String()
}

// strawberryTile is a strawberry: a low mat of leaves, white flowers, then
// red berries.
func strawberryTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		sprout(d, 5, 4, g.lean)
		for _, p := range [][2]float64{{-3, -1.6}, {2.8, -2}, {0.4, -3.4}} {
			if g.stage == 3 {
				blob{x + p[0], y + p[1], 1.4, 1.6, red, redShade}.draw(d)
			} else {
				blob{x + p[0], y + p[1], 1, 0.9, "#ffffff", "#f1e4f4"}.draw(d)
			}
		}
	}
	return d.String()
}
