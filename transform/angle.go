package transform

import "math"

// Angle is the transform's turn around the Z axis, in radians: the
// rotation a 2D game sees, clockwise on a y-down screen.
func (t Transform) Angle() float64 {
	q := t.Rotation
	return 2 * math.Atan2(q.Z, q.W)
}
