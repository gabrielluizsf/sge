package animal

import "math"

// Vec is a point or a direction.
type Vec [3]float64

// Quat is a rotation, x, y, z, w.
type Quat [4]float64

var identity = Quat{0, 0, 0, 1}

func add(a, b Vec) Vec           { return Vec{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func sub(a, b Vec) Vec           { return Vec{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func scale(a Vec, k float64) Vec { return Vec{a[0] * k, a[1] * k, a[2] * k} }
func dist(a, b Vec) float64 {
	d := sub(a, b)
	return math.Sqrt(d[0]*d[0] + d[1]*d[1] + d[2]*d[2])
}

// About is the turn by angle radians about axis, right-handed.
func About(axis Vec, angle float64) Quat {
	s, c := math.Sincos(angle / 2)
	n := math.Sqrt(axis[0]*axis[0] + axis[1]*axis[1] + axis[2]*axis[2])
	if n == 0 {
		return identity
	}
	return Quat{axis[0] / n * s, axis[1] / n * s, axis[2] / n * s, c}
}

// mul is the turn b followed by a.
func mul(a, b Quat) Quat {
	return Quat{
		a[3]*b[0] + a[0]*b[3] + a[1]*b[2] - a[2]*b[1],
		a[3]*b[1] - a[0]*b[2] + a[1]*b[3] + a[2]*b[0],
		a[3]*b[2] + a[0]*b[1] - a[1]*b[0] + a[2]*b[3],
		a[3]*b[3] - a[0]*b[0] - a[1]*b[1] - a[2]*b[2],
	}
}

func conj(q Quat) Quat { return Quat{-q[0], -q[1], -q[2], q[3]} }

// turn is v turned by q.
func turn(q Quat, v Vec) Vec {
	p := mul(mul(q, Quat{v[0], v[1], v[2], 0}), conj(q))
	return Vec{p[0], p[1], p[2]}
}

// sameSide flips q to the hemisphere of ref, so keyframes between them
// turn the short way.
func sameSide(q, ref Quat) Quat {
	if q[0]*ref[0]+q[1]*ref[1]+q[2]*ref[2]+q[3]*ref[3] < 0 {
		return Quat{-q[0], -q[1], -q[2], -q[3]}
	}
	return q
}

// wrap is angle a brought into -pi to pi.
func wrap(a float64) float64 {
	for a > math.Pi {
		a -= 2 * math.Pi
	}
	for a < -math.Pi {
		a += 2 * math.Pi
	}
	return a
}
