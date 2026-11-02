package art

// stalks draws n upright stalks of crop colour, tall high, leant, spread
// either side of the cell's middle.
func stalks(d *doc, n int, tall float64) {
	x, y := iso(Size/2, Size/2, 0)
	for i := range n {
		off := (float64(i) - float64(n-1)/2) * 2.2
		rod{x + off, y - 0.4, off * 0.1, -tall * (0.85 + 0.15*float64(i%2)), 1, "#8bc65a"}.draw(d)
	}
}

// cornTile is corn: a stalk of arching leaves, then a tassel and ears.
func cornTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		tall := 8 + 8*float64(g.stage-2)
		rod{x, y - 0.4, g.lean * 0.5, -tall, 1.2, leafBase}.draw(d)
		for i := 1.0; i < tall; i += 3.4 {
			side := 1.0 - 2*float64(int(i/3.4)%2)
			leaf(d, x+g.lean*0.5*i/tall, y-i, [2]float64{5 * side, -2})
		}
		if g.stage == 3 {
			rod{x + g.lean*0.5, y - tall, 0, -2.4, 0.6, "#d8b85a"}.draw(d)
			blob{x + 2.4, y - tall*0.55, 1.2, 2.6, yellow, "#d6a428"}.draw(d)
		}
	}
	return d.String()
}

// wheatTile is wheat: grass-green blades, then golden stalks with heads.
func wheatTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	base(d)
	switch {
	case seedling(d, g):
	case g.stage == 2:
		stalks(d, 5, 8)
	case g.stage == 3:
		x, y := iso(Size/2, Size/2, 0)
		for i := range 5 {
			off := (float64(i) - 2) * 2.2
			rod{x + off + g.lean*0.4, y - 0.4, 0, -12, 0.8, "#e6c46a"}.draw(d)
			blob{x + off + g.lean*0.4, y - 13, 0.9, 2.2, "#f4dc88", "#d4aa48"}.draw(d)
		}
	}
	return d.String()
}

// sugarcaneTile is sugar cane: a few jointed canes with a plume of leaves.
func sugarcaneTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		tall := 12 + 10*float64(g.stage-2)
		for i := range 3 {
			off := float64(i-1) * 2.6
			rod{x + off, y - 0.4, g.lean * 0.3, -tall, 1.5, "#a8d070"}.draw(d)
			for j := 3.0; j < tall; j += 4 {
				d.add(`<path d="M%g %g l1.6 0" stroke="%s" stroke-width="0.4"/>`, x+off-0.8, y-j, leafShade)
			}
			leaf(d, x+off+g.lean*0.3, y-tall, [2]float64{-4 + off, -3})
			leaf(d, x+off+g.lean*0.3, y-tall, [2]float64{4 + off, -3})
		}
	}
	return d.String()
}

// sunflowerTile is a sunflower: a tall stem with broad leaves, a bud, then
// a big yellow head.
func sunflowerTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		tall := 10 + 8*float64(g.stage-2)
		rod{x, y - 0.4, g.lean * 0.5, -tall, 1.1, leafBase}.draw(d)
		leaf(d, x, y-4, [2]float64{-4.5, -2})
		leaf(d, x, y-7, [2]float64{4.5, -2})
		hx, hy := x+g.lean*0.5, y-tall
		if g.stage == 2 {
			blob{hx, hy, 1.8, 1.8, leafBase, leafShade}.draw(d)
			break
		}
		for k := range 8 {
			a := float64(k) * 0.785
			blob{hx + 3.6*cosOf(a), hy + 3.6*sinOf(a), 1.5, 1.1, yellow, "#e0a820"}.draw(d)
		}
		blob{hx, hy, 2.4, 2.4, "#7a4b2a", "#4a2c1a"}.draw(d)
	}
	return d.String()
}

// cassavaTile is cassava: a reddish stem with fingered leaves, its root
// showing at the soil when ripe.
func cassavaTile(g growth) string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	switch {
	case seedling(d, g):
	case g.stage >= 2:
		tall := 8 + 6*float64(g.stage-2)
		rod{x, y - 0.4, g.lean * 0.5, -tall, 1, "#b0645a"}.draw(d)
		for k := range 5 {
			leaf(d, x+g.lean*0.5, y-tall, [2]float64{float64(k-2) * 2.4, -4 + float64(k%2)})
		}
		if g.stage == 3 {
			blob{x + 3, y - 0.4, 2.8, 1.2, "#c8a070", "#8a6540"}.draw(d)
		}
	}
	return d.String()
}
