package bake

import "math"

// fill calls plot for every pixel of a w by h grid whose middle lies in
// the triangle p, its corners' x and y in pixels, with the pixel's
// barycentric weights.
func fill(p [3][3]float64, w, h int, plot func(x, y int, wt [3]float64)) {
	x0 := max(0, int(math.Floor(min(p[0][0], p[1][0], p[2][0]))))
	x1 := min(w-1, int(math.Ceil(max(p[0][0], p[1][0], p[2][0]))))
	y0 := max(0, int(math.Floor(min(p[0][1], p[1][1], p[2][1]))))
	y1 := min(h-1, int(math.Ceil(max(p[0][1], p[1][1], p[2][1]))))
	area := edge(p[0], p[1], p[2])
	if area == 0 {
		return
	}
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			c := [3]float64{float64(x) + 0.5, float64(y) + 0.5}
			w0, w1, w2 := edge(p[1], p[2], c)/area, edge(p[2], p[0], c)/area, edge(p[0], p[1], c)/area
			if w0 >= 0 && w1 >= 0 && w2 >= 0 {
				plot(x, y, [3]float64{w0, w1, w2})
			}
		}
	}
}

// edge is twice the signed area of a, b, c.
func edge(a, b, c [3]float64) float64 {
	return (b[0]-a[0])*(c[1]-a[1]) - (b[1]-a[1])*(c[0]-a[0])
}
