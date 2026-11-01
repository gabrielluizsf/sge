package tree

import "math"

// vec is a point or direction in the tree's space: metres, Y up, its
// foot on the origin.
type vec = [3]float64

func add(a, b vec) vec           { return vec{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func sub(a, b vec) vec           { return vec{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func scale(a vec, k float64) vec { return vec{a[0] * k, a[1] * k, a[2] * k} }
func dot(a, b vec) float64       { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
func length(a vec) float64       { return math.Sqrt(dot(a, a)) }
func f32(a vec) [3]float32       { return [3]float32{float32(a[0]), float32(a[1]), float32(a[2])} }

// up is the way trees grow.
var up = vec{0, 1, 0}

func normalize(a vec) vec {
	if l := length(a); l > 0 {
		return scale(a, 1/l)
	}
	return a
}

func cross(a, b vec) vec {
	return vec{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

// perp is a unit vector square to a.
func perp(a vec) vec {
	u := normalize(cross(a, vec{1, 0, 0}))
	if length(u) < 0.5 {
		u = normalize(cross(a, vec{0, 0, 1}))
	}
	return u
}

// cone is a direction leaving axis by angle radians, turned round it to
// azimuth.
type cone struct {
	axis           vec
	angle, azimuth float64
}

// leave is the direction c says.
func leave(c cone) vec {
	axis := normalize(c.axis)
	u := perp(axis)
	w := cross(axis, u)
	side := add(scale(u, math.Cos(c.azimuth)), scale(w, math.Sin(c.azimuth)))
	return add(scale(axis, math.Cos(c.angle)), scale(side, math.Sin(c.angle)))
}
