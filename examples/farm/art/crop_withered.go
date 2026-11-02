package art

// witheredTile is a plant that went a day without water: a bent, brown
// stem and curled, dry leaves over the cracked mound.
func witheredTile() string {
	d := newDoc(IsoW, IsoH)
	x, y := base(d)
	rod{x, y - 0.4, 2.4, -4, 0.8, "#9a7a48"}.draw(d)
	rod{x + 2.4, y - 4.4, 2, 1.8, 0.7, "#9a7a48"}.draw(d)
	for _, t := range [][2]float64{{-3.4, 1.2}, {4.4, 3}, {-2, -1}} {
		d.add(`<path d="M%g %g q%g %g %g %g q-1 0.6 -2.2 0.2 Z" fill="#b89a5a" stroke="%s" stroke-width="0.4" stroke-linejoin="round"/>`, x+0.4, y-1, t[0]/2, -1.4, t[0], t[1]-2, line)
	}
	d.add(`<path d="M%g %g l1.4 0.6 M%g %g l-1.2 0.8" stroke="%s" stroke-width="0.35"/>`, x-2.6, y+0.3, x+2.8, y+0.1, soilShade)
	return d.String()
}
