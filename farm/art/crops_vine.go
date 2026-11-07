package art

// vined draws the vines and broad leaves a ground crop spreads, the
// further out the later its stage.
func vined(d *doc, g growth) {
	x, y := iso(Size/2, Size/2, 0)
	reach := 2.5 + 1.6*float64(g.stage)
	for _, s := range []float64{-1, 1} {
		rod{x, y - 0.4, s * reach, -0.8, 0.6, leafShade}.draw(d)
		leaf(d, x+s*reach*0.5, y-1, [2]float64{s * 3.4, -3.4})
		leaf(d, x+s*reach, y-1, [2]float64{s * 3, -2.4})
	}
	leaf(d, x, y-0.4, [2]float64{0, -4.4})
}

// pumpkinTile is a pumpkin: vines and leaves, a small green fruit, then a
// big orange ribbed one.
func pumpkinTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		vined(d, g)
		fill, shade := "#8cb84a", "#5a8a32"
		r := 2.0
		if g.stage == 3 {
			fill, shade, r = orange, orangeShade, 3.6
		}
		blob{x + 1, y - r*0.5, r, r * 0.8, fill, shade}.draw(d)
		d.add(`<path d="M%g %g q-0.8 1.6 0 3.2 M%g %g q0.8 1.6 0 3.2" stroke="%s" stroke-width="0.35" fill="none"/>`, x+0.2, y-r*1.2, x+1.8, y-r*1.2, shade)
		rod{x + 1, y - r*1.3, 0.4, -1.2, 0.7, leafShade}.draw(d)
	}
	return d.String()
}

// watermelonTile is a watermelon: vines, a small round melon, then a big
// striped one.
func watermelonTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		vined(d, g)
		r := 2.0
		if g.stage == 3 {
			r = 4
		}
		blob{x + 1, y - r*0.5, r * 1.2, r * 0.85, "#5fb14a", "#2f7a3a"}.draw(d)
		for _, s := range []float64{-0.5, 0.4, 1.3} {
			d.add(`<path d="M%g %g q0.6 %g 0 %g" stroke="%s" stroke-width="0.5" fill="none"/>`, x+1+s*r*0.6, y-r*1.1, r*0.5, r*1.1, "#245f2c")
		}
	}
	return d.String()
}
