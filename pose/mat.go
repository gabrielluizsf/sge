// Package pose works out where a glTF model's nodes are at a moment of one
// of its animations, and the matrices its skins' joints bend their
// vertices by: the arithmetic a renderer does for every model every frame,
// kept apart from any one renderer.
package pose

// Mat is a 4 by 4 matrix, column-major, as glTF writes them.
type Mat [16]float64

// Identity is the matrix that moves nothing.
var Identity = Mat{0: 1, 5: 1, 10: 1, 15: 1}

// Mul is a times b: b's change, then a's.
func (a Mat) Mul(b Mat) Mat {
	var m Mat
	for c := range 4 {
		for r := range 4 {
			m[c*4+r] = a[r]*b[c*4] + a[4+r]*b[c*4+1] + a[8+r]*b[c*4+2] + a[12+r]*b[c*4+3]
		}
	}
	return m
}

// Affine is a times b for matrices whose last row is 0, 0, 0, 1, as every
// node's matrix and every inverse bind matrix of glTF is: a quarter of
// Mul's work left out.
func (a Mat) Affine(b Mat) Mat {
	var m Mat
	MulAffine(&m, &a, &b)
	return m
}

// MulAffine sets dst to a times b, as Affine does, without copying them;
// dst may be a or b.
func MulAffine(dst, a, b *Mat) {
	// a is read whole first, and each column of b before dst's is
	// written, so that dst may be either.
	a0, a1, a2, a4, a5, a6, a8, a9, a10 := a[0], a[1], a[2], a[4], a[5], a[6], a[8], a[9], a[10]
	a12, a13, a14 := a[12], a[13], a[14]
	x, y, z := b[0], b[1], b[2]
	dst[0], dst[1], dst[2], dst[3] = a0*x+a4*y+a8*z, a1*x+a5*y+a9*z, a2*x+a6*y+a10*z, 0
	x, y, z = b[4], b[5], b[6]
	dst[4], dst[5], dst[6], dst[7] = a0*x+a4*y+a8*z, a1*x+a5*y+a9*z, a2*x+a6*y+a10*z, 0
	x, y, z = b[8], b[9], b[10]
	dst[8], dst[9], dst[10], dst[11] = a0*x+a4*y+a8*z, a1*x+a5*y+a9*z, a2*x+a6*y+a10*z, 0
	x, y, z = b[12], b[13], b[14]
	dst[12], dst[13], dst[14], dst[15] = a0*x+a4*y+a8*z+a12, a1*x+a5*y+a9*z+a13, a2*x+a6*y+a10*z+a14, 1
}

// Point is p moved by a.
func (a Mat) Point(p [3]float64) [3]float64 {
	return [3]float64{
		a[0]*p[0] + a[4]*p[1] + a[8]*p[2] + a[12],
		a[1]*p[0] + a[5]*p[1] + a[9]*p[2] + a[13],
		a[2]*p[0] + a[6]*p[1] + a[10]*p[2] + a[14],
	}
}

// Dir is the direction d turned by a, not moved.
func (a Mat) Dir(d [3]float64) [3]float64 {
	return [3]float64{
		a[0]*d[0] + a[4]*d[1] + a[8]*d[2],
		a[1]*d[0] + a[5]*d[1] + a[9]*d[2],
		a[2]*d[0] + a[6]*d[1] + a[10]*d[2],
	}
}

// TRS is the matrix of a translation, a rotation (a quaternion x, y, z,
// w) and a scale, applied scale first.
func TRS(t [3]float64, q [4]float64, s [3]float64) Mat {
	x, y, z, w := q[0], q[1], q[2], q[3]
	return Mat{
		(1 - 2*(y*y+z*z)) * s[0], 2 * (x*y + z*w) * s[0], 2 * (x*z - y*w) * s[0], 0,
		2 * (x*y - z*w) * s[1], (1 - 2*(x*x+z*z)) * s[1], 2 * (y*z + x*w) * s[1], 0,
		2 * (x*z + y*w) * s[2], 2 * (y*z - x*w) * s[2], (1 - 2*(x*x+y*y)) * s[2], 0,
		t[0], t[1], t[2], 1,
	}
}
