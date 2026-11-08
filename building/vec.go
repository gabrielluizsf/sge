package building

import "math"

// vec is a point or direction in a building's space: metres, Y up, X
// along its front, Z out of its front.
type vec = [3]float64

// up is the way up.
var up = vec{0, 1, 0}

func add(a, b vec) vec           { return vec{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func sub(a, b vec) vec           { return vec{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func scale(a vec, k float64) vec { return vec{a[0] * k, a[1] * k, a[2] * k} }
func dot(a, b vec) float64       { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }
func length(a vec) float64       { return math.Sqrt(dot(a, a)) }
func f32(a vec) [3]float32       { return [3]float32{float32(a[0]), float32(a[1]), float32(a[2])} }

func cross(a, b vec) vec {
	return vec{a[1]*b[2] - a[2]*b[1], a[2]*b[0] - a[0]*b[2], a[0]*b[1] - a[1]*b[0]}
}

func normalize(a vec) vec {
	if l := length(a); l > 0 {
		return scale(a, 1/l)
	}
	return a
}

// along is the point k of the way from a to b.
func along(a, b vec, k float64) vec { return add(a, scale(sub(b, a), k)) }

// centroid is the average of points.
func centroid(points []vec) vec {
	var c vec
	for _, p := range points {
		c = add(c, p)
	}
	return scale(c, 1/float64(max(len(points), 1)))
}

// newell is the normal of a flat polygon, its corners counterclockwise
// seen from the side it faces, by Newell's method, which a polygon of
// more than three corners needs; it is of length 1, or 0 for no area.
func newell(corners []vec) vec {
	var n vec
	for i, a := range corners {
		b := corners[(i+1)%len(corners)]
		n[0] += (a[1] - b[1]) * (a[2] + b[2])
		n[1] += (a[2] - b[2]) * (a[0] + b[0])
		n[2] += (a[0] - b[0]) * (a[1] + b[1])
	}
	return normalize(n)
}
