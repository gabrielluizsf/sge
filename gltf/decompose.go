package gltf

import "math"

// decompose is the translation, rotation and scale of m, a column-major
// matrix of a translation, rotation and positive scale, as glTF's are.
func decompose(m [16]float64) (t [3]float64, r [4]float64, s [3]float64) {
	t = [3]float64{m[12], m[13], m[14]}
	col := func(c int) [3]float64 { return [3]float64{m[c*4], m[c*4+1], m[c*4+2]} }
	var rot [3][3]float64
	for c := range 3 {
		v := col(c)
		s[c] = math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
		for k := range 3 {
			if s[c] != 0 {
				rot[c][k] = v[k] / s[c]
			}
		}
	}
	// rot[c][k] is row k of column c; the quaternion of that rotation.
	m00, m11, m22 := rot[0][0], rot[1][1], rot[2][2]
	switch tr := m00 + m11 + m22; {
	case tr > 0:
		q := math.Sqrt(tr+1) * 2
		r = [4]float64{(rot[1][2] - rot[2][1]) / q, (rot[2][0] - rot[0][2]) / q, (rot[0][1] - rot[1][0]) / q, q / 4}
	case m00 > m11 && m00 > m22:
		q := math.Sqrt(1+m00-m11-m22) * 2
		r = [4]float64{q / 4, (rot[1][0] + rot[0][1]) / q, (rot[2][0] + rot[0][2]) / q, (rot[1][2] - rot[2][1]) / q}
	case m11 > m22:
		q := math.Sqrt(1+m11-m00-m22) * 2
		r = [4]float64{(rot[1][0] + rot[0][1]) / q, q / 4, (rot[2][1] + rot[1][2]) / q, (rot[2][0] - rot[0][2]) / q}
	default:
		q := math.Sqrt(1+m22-m00-m11) * 2
		r = [4]float64{(rot[2][0] + rot[0][2]) / q, (rot[2][1] + rot[1][2]) / q, q / 4, (rot[0][1] - rot[1][0]) / q}
	}
	return t, r, s
}
