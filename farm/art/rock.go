package art

// RockVariants is how many rocks there are: a small rock, a boulder and a
// few rocks together.
const RockVariants = 3

// rockTile is rock v standing on its cell: grey stone, light on top and
// shaded below, with a crack or two.
func rockTile(v int) string {
	d := newDoc(IsoW, IsoH)
	x, y := iso(Size/2, Size/2, 0)
	switch v {
	case 1:
		shadow(d, 8, 3.4)
		stone(d, x, y, 1.45)
	case 2:
		shadow(d, 9, 3.6)
		stone(d, x-4.4, y+1, 0.8)
		stone(d, x+3.8, y+0.6, 0.95)
		stone(d, x-0.4, y-0.8, 0.7)
	default:
		shadow(d, 5, 2.2)
		stone(d, x, y, 0.85)
	}
	return d.String()
}

// stone is one rock of scale s standing at x, y: a rounded lump, its top
// lit, with a crack.
func stone(d *doc, x, y, s float64) {
	body := d.down("#cfcac2", "#77716a")
	d.add(`<path d="M%.2f %.2f q%.2f %.2f %.2f %.2f q%.2f %.2f %.2f %.2f q%.2f %.2f %.2f %.2f q%.2f %.2f %.2f %.2f Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`,
		x-6*s, y, -0.6*s, -6*s, 3*s, -7.2*s, 4*s, -2.6*s, 7.6*s, 0.4*s, 2.2*s, 2.8*s, 1.4*s, 6.8*s, -7*s, 2.2*s, -9*s, -2.2*s, body, line)
	d.add(`<path d="M%.2f %.2f q%.2f %.2f %.2f %.2f" stroke="%s" stroke-width="0.4" fill="none" stroke-linecap="round"/>`, x+0.5*s, y-5*s, 1.2*s, 1.8*s, -0.4*s, 3.6*s, line)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="%.2f" ry="%.2f" fill="%s" opacity="0.55"/>`, x-1.8*s, y-5.4*s, 1.6*s, 0.8*s, white)
}
