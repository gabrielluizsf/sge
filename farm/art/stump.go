package art

// stumpTile is what is left of a tree cut down, standing on its cell: a
// stump with its rings showing, and chips about it.
func stumpTile() string {
	d := newDoc(IsoW, IsoH)
	shadow(d, 5, 2.4)
	x, y := iso(Size/2, Size/2, 0)
	bark := d.linear(grad{0, 0, 1, 0, []stop{{0, barkLight}, {0.55, barkBase}, {1, barkShade}}})
	d.add(`<path d="M%.2f %.2f l0 -4 a3.4 1.6 0 0 1 6.8 0 l0 4 a3.4 1.6 0 0 1 -6.8 0 Z" fill="%s" stroke="%s" stroke-width="0.6"/>`, x-3.4, y, bark, line)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="3.4" ry="1.6" fill="%s" stroke="%s" stroke-width="0.5"/>`, x, y-4, woodLight, line)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="1.9" ry="0.85" fill="none" stroke="%s" stroke-width="0.35"/>`, x, y-4, woodShade)
	d.add(`<circle cx="%.2f" cy="%.2f" r="0.35" fill="%s"/>`, x, y-4, woodShade)
	for _, c := range [][2]float64{{-5, 1}, {4.6, 0.4}, {2, 2.2}} {
		d.add(`<path d="M%.2f %.2f l1.2 -0.4 l0.4 0.6 Z" fill="%s" stroke="%s" stroke-width="0.25"/>`, x+c[0], y+c[1], woodLight, line)
	}
	return d.String()
}
