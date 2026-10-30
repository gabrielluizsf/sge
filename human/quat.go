package human

import "math"

// quat is a rotation, x, y, z and w, as glTF writes them.
type quat [4]float64

var still = quat{0, 0, 0, 1}

// around is a turn of angle about axis, of length 1.
func around(axis [3]float64, angle float64) quat {
	s := math.Sin(angle / 2)
	return quat{axis[0] * s, axis[1] * s, axis[2] * s, math.Cos(angle / 2)}
}

// then is q followed by r: r turning what q has turned.
func (q quat) then(r quat) quat {
	return quat{
		r[3]*q[0] + r[0]*q[3] + r[1]*q[2] - r[2]*q[1],
		r[3]*q[1] - r[0]*q[2] + r[1]*q[3] + r[2]*q[0],
		r[3]*q[2] + r[0]*q[1] - r[1]*q[0] + r[2]*q[3],
		r[3]*q[3] - r[0]*q[0] - r[1]*q[1] - r[2]*q[2],
	}
}

// between is the least turn taking direction a to b.
func between(a, b [3]float64) quat {
	a, b = normalize(a), normalize(b)
	c := cross(a, b)
	d := a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
	if d < -0.999999 {
		return around(normalize(cross(a, [3]float64{0, 0, 1})), math.Pi)
	}
	q := quat{c[0], c[1], c[2], 1 + d}
	l := math.Sqrt(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])
	return quat{q[0] / l, q[1] / l, q[2] / l, q[3] / l}
}

// turn is v turned by q.
func (q quat) turn(v [3]float64) [3]float64 {
	u := [3]float64{q[0], q[1], q[2]}
	t := scale(cross(u, v), 2)
	return add(add(v, scale(t, q[3])), cross(u, t))
}

// dot4 is a key's rotation dotted with q: below 0, they are on opposite
// sides, and a key between them would go the long way round.
func dot4(k [4]float32, q quat) float64 {
	return float64(k[0])*q[0] + float64(k[1])*q[1] + float64(k[2])*q[2] + float64(k[3])*q[3]
}
