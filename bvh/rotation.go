package bvh

import "math"

// Rotation is joint i's turn in frame f, x, y, z and w, as glTF writes a
// rotation, in the file's own axes: its rotation channels composed in the
// order they come, the first outermost, as BVH composes them.
func (m *Motion) Rotation(f, i int) [4]float64 {
	q := [4]float64{0, 0, 0, 1}
	for k, v := range m.values(f, i) {
		axis := -1
		switch m.Joints[i].Channels[k] {
		case Xrotation:
			axis = 0
		case Yrotation:
			axis = 1
		case Zrotation:
			axis = 2
		}
		if axis < 0 {
			continue
		}
		half := v * math.Pi / 360
		var r [4]float64
		r[axis], r[3] = math.Sin(half), math.Cos(half)
		q = mul(q, r)
	}
	return q
}

// mul is the Hamilton product a b: b turned first, then a.
func mul(a, b [4]float64) [4]float64 {
	return [4]float64{
		a[3]*b[0] + a[0]*b[3] + a[1]*b[2] - a[2]*b[1],
		a[3]*b[1] - a[0]*b[2] + a[1]*b[3] + a[2]*b[0],
		a[3]*b[2] + a[0]*b[1] - a[1]*b[0] + a[2]*b[3],
		a[3]*b[3] - a[0]*b[0] - a[1]*b[1] - a[2]*b[2],
	}
}
