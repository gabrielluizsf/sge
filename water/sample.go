package water

import (
	"image"
	"math"
)

// normalAt is the normal of map m at (u, v), the map taken to repeat
// every 1, the four pixels about it blended; unit length, X along u, Y
// along v, Z up.
func normalAt(m *image.NRGBA, u, v float64) [3]float64 {
	w, h := m.Bounds().Dx(), m.Bounds().Dy()
	x, y := u*float64(w)-0.5, v*float64(h)-0.5
	x0, y0 := math.Floor(x), math.Floor(y)
	fx, fy := x-x0, y-y0
	var n [3]float64
	for i, wt := range [4]float64{(1 - fx) * (1 - fy), fx * (1 - fy), (1 - fx) * fy, fx * fy} {
		px := (int(x0) + i%2 + w*4) % w
		py := (int(y0) + i/2 + h*4) % h
		c := m.NRGBAAt(px, py)
		n[0] += wt * (float64(c.R)/255*2 - 1)
		n[1] += wt * (float64(c.G)/255*2 - 1)
		n[2] += wt * (float64(c.B)/255*2 - 1)
	}
	return unit(n)
}

// unit is n with length 1.
func unit(n [3]float64) [3]float64 {
	l := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	if l == 0 {
		return [3]float64{0, 0, 1}
	}
	return [3]float64{n[0] / l, n[1] / l, n[2] / l}
}

func dot(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
