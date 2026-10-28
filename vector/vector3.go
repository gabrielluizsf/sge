// Package vector holds the 3D math the engine moves things with: points and
// directions ([Vector3]) and rotations ([Quaternion]).
package vector

// Vector3 is a point or a direction in 3D space.
type Vector3 struct {
	X, Y, Z float64
}
