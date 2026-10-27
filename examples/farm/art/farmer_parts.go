package art

// The farmer's drawing, after the farm's character sheet, FarmerW by
// FarmerH units, standing on the bottom edge: a straw hat with a brown
// band over brown hair, a cream shirt with its sleeves rolled up, blue
// denim overalls with brass buttons and boots of brown leather, outlined
// and shaded as anime is.
const (
	FarmerW = 24
	FarmerH = 40
)

// pt is a point of the farmer's drawing.
type pt struct{ x, y float64 }

// limb is a rounded stroke from a to b, w wide, in c, outlined.
type limb struct {
	a, b pt
	w    float64
	c    string
}

// draw draws l: the outline a little wider, under the colour.
func (l limb) draw(d *doc) {
	for _, s := range []struct {
		w float64
		c string
	}{{l.w + 1.1, line}, {l.w, l.c}} {
		d.add(`<path d="M%.2f %.2f L%.2f %.2f" stroke="%s" stroke-width="%.2f" stroke-linecap="round"/>`, l.a.x, l.a.y, l.b.x, l.b.y, s.c, s.w)
	}
}

// boot is a boot whose sole's middle is at p, its toe pointing toe, 1 to
// the right, -1 to the left, 0 at the viewer.
func boot(d *doc, p pt, toe float64) {
	w := 2.2 + 1.2*abs(toe)
	cx := p.x + toe*1
	d.add(`<path d="M%.2f %.2f q0 -3 %.2f -3.2 l%.2f 0 q%.2f 0.2 %.2f 3.2 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`,
		cx-w, p.y, w*0.4, w*1.2, w*0.4, w*0.4, d.down(bootLight, bootBase), line)
	d.add(`<path d="M%.2f %.2f l%.2f 0" stroke="%s" stroke-width="0.8"/>`, cx-w+0.3, p.y-0.4, 2*w-0.6, bootShade)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// hat is the straw hat, its crown's middle at x, the brim at y; behind
// has the brim drawn whole, as from behind the crown hides none of it.
func hat(d *doc, x, y float64) {
	straw := d.linear(grad{0, 0, 0, 1, []stop{{0, strawLight}, {0.6, strawBase}, {1, strawShade}}})
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="9.2" ry="2.6" fill="%s" stroke="%s" stroke-width="0.7"/>`, x, y+0.6, straw, line)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="7.2" ry="1.4" fill="none" stroke="%s" stroke-width="0.35" opacity="0.7"/>`, x, y+0.6, strawShade)
	d.add(`<path d="M%.2f %.2f q0 -5.4 4.4 -5.4 q4.4 0 4.4 5.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, x-4.4, y, straw, line)
	d.add(`<path d="M%.2f %.2f q4.4 1.2 8.8 0 l0 -1.5 q-4.4 1.1 -8.8 0 Z" fill="%s"/>`, x-4.4, y, band)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="1.4" ry="0.9" fill="%s" opacity="0.7"/>`, x-1.8, y-3.6, white)
}

// eyeAt is an anime eye at p: tall and dark, a light in its upper side.
func eyeAt(d *doc, p pt) {
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="1.15" ry="1.6" fill="%s"/>`, p.x, p.y, eye)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="0.7" ry="0.75" fill="%s" opacity="0.9"/>`, p.x, p.y+0.6, hairLight)
	d.add(`<circle cx="%.2f" cy="%.2f" r="0.45" fill="%s"/>`, p.x-0.35, p.y-0.55, eyeLight)
	d.add(`<circle cx="%.2f" cy="%.2f" r="0.2" fill="%s"/>`, p.x+0.4, p.y+0.7, eyeLight)
	d.add(`<path d="M%.2f %.2f q1.25 -0.8 2.5 0" stroke="%s" stroke-width="0.55" fill="none" stroke-linecap="round"/>`, p.x-1.25, p.y-1.6, line)
}

// faceAt is a face whose middle is at c, looking at the viewer.
func faceAt(d *doc, c pt) {
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="4.3" ry="4.7" fill="%s" stroke="%s" stroke-width="0.7"/>`, c.x, c.y, skin, line)
	d.add(`<path d="M%.2f %.2f q4.3 3.2 8.6 0 q-0.4 3.6 -4.3 4.6 q-3.9 -1 -4.3 -4.6 Z" fill="%s" opacity="0.5"/>`, c.x-4.3, c.y+0.2, skinShade)
	eyeAt(d, pt{c.x - 1.9, c.y + 1})
	eyeAt(d, pt{c.x + 1.9, c.y + 1})
	for _, s := range []float64{-1, 1} {
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="0.9" ry="0.45" fill="%s" opacity="0.55"/>`, c.x+s*2.9, c.y+2.9, blush)
	}
	d.add(`<path d="M%.2f %.2f q0.6 0.45 1.2 0" stroke="%s" stroke-width="0.4" fill="none" stroke-linecap="round"/>`, c.x-0.6, c.y+3.6, line)
}

// bangs is brown hair over the brow, spiky, from x0 to x1 at the brow's
// line y.
func bangs(d *doc, x0, x1, y float64) {
	d.add(`<path d="M%.2f %.2f q%.2f -3.6 %.2f -0.4 l-0.8 1.6 l-1 -1.4 l-0.9 1.9 l-1.1 -1.8 l-0.9 1.6 l-1.1 -1.7 l-0.6 1.8 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`,
		x0, y+1.6, (x1-x0)/2, x1-x0, d.down(hairLight, hairBase), line)
}
