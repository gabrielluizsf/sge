package art

// growth is how far a crop is drawn: its stage, from 0 the seed to 3
// ripe, and how far the wind leans it.
type growth struct {
	stage int
	lean  float64
}

// blob is a rounded thing on a plant, a fruit, a bud or a head, at x, y,
// rx by ry, in a colour and its shade.
type blob struct {
	x, y, rx, ry float64
	fill, shade  string
}

// draw paints the blob, outlined, with a glint.
func (b blob) draw(d *doc) {
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" stroke="%s" stroke-width="0.4"/>`, b.x, b.y, b.rx, b.ry, d.down(b.fill, b.shade), line)
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" opacity="0.6"/>`, b.x-b.rx*0.35, b.y-b.ry*0.35, b.rx*0.28, b.ry*0.25, white)
}

// rod is a straight stalk from x, y by dx, dy, w wide, in a colour.
type rod struct {
	x, y, dx, dy float64
	w            float64
	colour       string
}

// draw paints the rod with round ends.
func (r rod) draw(d *doc) {
	d.add(`<path d="M%g %g l%g %g" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`, r.x, r.y, r.dx, r.dy, line, r.w+0.6)
	d.add(`<path d="M%g %g l%g %g" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`, r.x, r.y, r.dx, r.dy, r.colour, r.w)
}

// base starts a crop's tile: its shadow and the mound it grows from, and
// where the plant stands.
func base(d *doc) (x, y float64) {
	shadow(d, 4, 1.6)
	mound(d)
	return iso(Size/2, Size/2, 0)
}

// seedling draws the first stages that all crops share: stage 1 and, for
// low crops, a bigger tuft. It reports whether it drew the stage.
func seedling(d *doc, g growth) bool {
	if g.stage != 1 {
		return false
	}
	sprout(d, 2, 3, g.lean)
	return true
}
