package render3d

import "math"

// frustum is the six planes bounding what a camera sees, each a, b, c, d
// with a, b, c of length 1 and a point p inside when a·p + d ≥ 0.
type frustum [6][4]float64

// frustumOf is the frustum of c, from its ViewProj's rows: left, right,
// bottom, top, near (depth 0) and far (depth 1).
func frustumOf(c *Camera) frustum {
	m := &c.ViewProj
	row := func(i int) [4]float64 { return [4]float64{m[i], m[4+i], m[8+i], m[12+i]} }
	r0, r1, r2, r3 := row(0), row(1), row(2), row(3)
	var f frustum
	for k := range 4 {
		f[0][k] = r3[k] + r0[k]
		f[1][k] = r3[k] - r0[k]
		f[2][k] = r3[k] + r1[k]
		f[3][k] = r3[k] - r1[k]
		f[4][k] = r2[k]
		f[5][k] = r3[k] - r2[k]
	}
	for i := range f {
		l := math.Sqrt(f[i][0]*f[i][0] + f[i][1]*f[i][1] + f[i][2]*f[i][2])
		if l > 0 {
			f[i] = [4]float64{f[i][0] / l, f[i][1] / l, f[i][2] / l, f[i][3] / l}
		}
	}
	return f
}

// sees is whether any of s may be inside f.
func (f *frustum) sees(s Sphere) bool {
	for i := range f {
		p := &f[i]
		if p[0]*s.Center[0]+p[1]*s.Center[1]+p[2]*s.Center[2]+p[3] < -s.Radius {
			return false
		}
	}
	return true
}
