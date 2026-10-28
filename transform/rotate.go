package transform

import "github.com/gabrielluizsf/sge/vector"

// Rotate turns the transform a further angle radians around axis, on top of
// the rotation it already has.
func (t *Transform) Rotate(axis vector.Vector3, angle float64) {
	t.Rotation = vector.AxisAngle(axis, angle).Mul(t.Rotation)
}
