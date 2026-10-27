package art

import (
	"fmt"
	"strings"
)

// doc is an SVG drawing being written, w by h units: its gradients and
// its shapes, in the order drawn.
type doc struct {
	w, h       float64
	defs, body strings.Builder
	ids        int
}

func newDoc(w, h float64) *doc { return &doc{w: w, h: h} }

// add writes a shape, an SVG element fmt formats.
func (d *doc) add(format string, a ...any) {
	fmt.Fprintf(&d.body, format+"\n", a...)
}

// stop is a colour of a gradient at offset, from 0 to 1.
type stop struct {
	at    float64
	color string
}

// grad is a gradient from x1, y1 to x2, y2, as fractions of what it
// paints, through its stops.
type grad struct {
	x1, y1, x2, y2 float64
	stops          []stop
}

// linear defines g and is the paint that fills with it.
func (d *doc) linear(g grad) string {
	d.ids++
	id := fmt.Sprint("g", d.ids)
	fmt.Fprintf(&d.defs, `<linearGradient id="%s" x1="%g" y1="%g" x2="%g" y2="%g">`, id, g.x1, g.y1, g.x2, g.y2)
	for _, s := range g.stops {
		fmt.Fprintf(&d.defs, `<stop offset="%g" stop-color="%s"/>`, s.at, s.color)
	}
	d.defs.WriteString("</linearGradient>\n")
	return "url(#" + id + ")"
}

// down is a gradient from light at the top to dark at the bottom.
func (d *doc) down(light, dark string) string {
	return d.linear(grad{0, 0, 0, 1, []stop{{0, light}, {1, dark}}})
}

// String is the drawing as SVG.
func (d *doc) String() string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %g %g" width="%g" height="%g">`+"\n<defs>\n%s</defs>\n%s</svg>\n",
		d.w, d.h, d.w, d.h, d.defs.String(), d.body.String())
}

// onCell is the SVG transform drawing a cell's square, 16 by 16 from its
// far corner, as the diamond of a tile's picture: u across, v down.
const onCell = `matrix(1 0.5 -1 0.5 16 32)`

// iso is where the point u, v of a cell, h up from the ground, is drawn
// on a tile's picture.
func iso(u, v, h float64) (x, y float64) {
	return 16 + u - v, 32 + (u+v)/2 - h
}

// rng is a small deterministic random sequence, the same art every run.
type rng uint32

func (r *rng) next() float64 {
	*r = *r*1664525 + 1013904223
	return float64(*r>>8) / float64(1<<24)
}
