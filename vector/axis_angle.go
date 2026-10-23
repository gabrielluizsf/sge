package vector

import "math"

// AxisAngle is the rotation by angle radians around axis. The axis does not
// have to be normalized.
func AxisAngle(axis Vector3, angle float64) Quaternion {
	a := axis.Normalized()
	s := math.Sin(angle / 2)
	return Quaternion{a.X * s, a.Y * s, a.Z * s, math.Cos(angle / 2)}
}
