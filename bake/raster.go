package bake

import "math"

// fill calls plot for every pixel of a grid of size, width by height,
// whose middle lies in the triangle p, its corners' x and y in pixels,
// with the pixel's barycentric weights.
func fill(p [3][3]float64, size [2]int, plot func(x, y int, wt [3]float64)) {
	w, h := size[0], size[1]
	x0 := max(0, int(math.Floor(min(p[0][0], p[1][0], p[2][0]))))
	x1 := min(w-1, int(math.Ceil(max(p[0][0], p[1][0], p[2][0]))))
	y0 := max(0, int(math.Floor(min(p[0][1], p[1][1], p[2][1]))))
	y1 := min(h-1, int(math.Ceil(max(p[0][1], p[1][1], p[2][1]))))
	area := edge(p[0], p[1], p[2])
	if area == 0 {
		return
	}
	for y := y0; y <= y1; y++ {
		cy := float64(y) + 0.5
		for x := x0; x <= x1; x++ {
			cx := float64(x) + 0.5
			w0, w1, w2 := edgeAt(p[1], p[2], cx, cy)/area, edgeAt(p[2], p[0], cx, cy)/area, edgeAt(p[0], p[1], cx, cy)/area
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

// edgeAt is edge of a, b and the point cx, cy, the same sum in the same
// order, with no point made to hold it: fill asks it of every pixel.
func edgeAt(a, b [3]float64, cx, cy float64) float64 {
	return (b[0]-a[0])*(cy-a[1]) - (b[1]-a[1])*(cx-a[0])
}
