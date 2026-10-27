package art

import "math"

// raster draws on a grid in the cell's own terms: u across and v down
// the cell, h up from the ground, a cell 16 each way, its far corner at
// the grid's point ox, oy.
type raster struct {
	g      Grid
	ox, oy float64
}

// point is a point of the cell, and its height.
type point struct{ u, v, h float64 }

// at is where p is drawn on the grid.
func (r *raster) at(p point) (x, y float64) {
	return r.ox + p.u - p.v, r.oy + (p.u+p.v)/2 - p.h
}

// poly fills the polygon of the points with c, a pixel wherever the
// pixel's middle lies inside.
func (r *raster) poly(c rune, pts ...point) {
	xs, ys := make([]float64, len(pts)), make([]float64, len(pts))
	for i, p := range pts {
		xs[i], ys[i] = r.at(p)
	}
	for y := range r.g {
		for x := range len([]rune(r.g[y])) {
			if inside(xs, ys, float64(x)+0.5, float64(y)+0.5) {
				r.g = r.g.set(x, y, c)
			}
		}
	}
}

// inside reports whether x, y is inside the polygon, by the even-odd rule.
func inside(xs, ys []float64, x, y float64) bool {
	in := false
	for i, j := 0, len(xs)-1; i < len(xs); j, i = i, i+1 {
		if (ys[i] > y) != (ys[j] > y) && x < (xs[j]-xs[i])*(y-ys[i])/(ys[j]-ys[i])+xs[i] {
			in = !in
		}
	}
	return in
}

// line draws c from a to b, a pixel a step.
func (r *raster) line(c rune, a, b point) {
	x0, y0 := r.at(a)
	x1, y1 := r.at(b)
	n := int(math.Max(math.Abs(x1-x0), math.Abs(y1-y0))) + 1
	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		x, y := int(math.Floor(x0+(x1-x0)*t)), int(math.Floor(y0+(y1-y0)*t))
		if y >= 0 && y < len(r.g) && x >= 0 && x < len([]rune(r.g[y])) {
			r.g = r.g.set(x, y, c)
		}
	}
}
