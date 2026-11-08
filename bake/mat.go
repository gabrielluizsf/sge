package bake

import "math"

// mat is a 4 by 4 matrix, column-major, as glTF writes them.
type mat [16]float64

var identity = mat{0: 1, 5: 1, 10: 1, 15: 1}

// mul is a times b.
func (a mat) mul(b mat) mat {
	var m mat
	for c := range 4 {
		for r := range 4 {
			var s float64
			for k := range 4 {
				s += a[k*4+r] * b[c*4+k]
			}
			m[c*4+r] = s
		}
	}
	return m
}

// point is p moved by m; dir is a direction turned by m, not moved.
func (a mat) point(p [3]float64) [3]float64 {
	return [3]float64{
		a[0]*p[0] + a[4]*p[1] + a[8]*p[2] + a[12],
		a[1]*p[0] + a[5]*p[1] + a[9]*p[2] + a[13],
		a[2]*p[0] + a[6]*p[1] + a[10]*p[2] + a[14],
	}
}

func (a mat) dir(p [3]float64) [3]float64 {
	return [3]float64{
		a[0]*p[0] + a[4]*p[1] + a[8]*p[2],
		a[1]*p[0] + a[5]*p[1] + a[9]*p[2],
		a[2]*p[0] + a[6]*p[1] + a[10]*p[2],
	}
}

// trs is the matrix of a translation, rotation (x, y, z, w) and scale.
func trs(t [3]float64, q [4]float64, s [3]float64) mat {
	x, y, z, w := q[0], q[1], q[2], q[3]
	return mat{
		(1 - 2*(y*y+z*z)) * s[0], 2 * (x*y + z*w) * s[0], 2 * (x*z - y*w) * s[0], 0,
		2 * (x*y - z*w) * s[1], (1 - 2*(x*x+z*z)) * s[1], 2 * (y*z + x*w) * s[1], 0,
		2 * (x*z + y*w) * s[2], 2 * (y*z - x*w) * s[2], (1 - 2*(x*x+y*y)) * s[2], 0,
		t[0], t[1], t[2], 1,
	}
}

// slerp is the rotation f of the way from a to b, the short way round.
func slerp(a, b [4]float64, f float64) [4]float64 {
	d := a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
	if d < 0 {
		b, d = [4]float64{-b[0], -b[1], -b[2], -b[3]}, -d
	}
	wa, wb := 1-f, f
	if d < 0.9995 {
		th := math.Acos(d)
		wa, wb = math.Sin((1-f)*th)/math.Sin(th), math.Sin(f*th)/math.Sin(th)
	}
	var q [4]float64
	var n float64
	for i := range q {
		q[i] = wa*a[i] + wb*b[i]
		n += q[i] * q[i]
	}
	n = math.Sqrt(n)
	for i := range q {
		q[i] /= n
	}
	return q
}

func normalize(v [3]float64) [3]float64 {
	l := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	if l == 0 {
		return v
	}
	return [3]float64{v[0] / l, v[1] / l, v[2] / l}
}

func dot(a, b [3]float64) float64 { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }

func cross(a, b [3]float64) [3]float64 {
	return [3]float64{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}
