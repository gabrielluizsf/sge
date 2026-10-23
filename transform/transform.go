// Package transform places a game object in the world: where it is, how it is
// turned and how big it is.
package transform

import "github.com/gabrielluizsf/sge/vector"

// Transform is an object's position, rotation and scale. The zero value is
// collapsed to a point; use [New] for an object at its natural size.
type Transform struct {
	Position vector.Vector3
	Rotation vector.Quaternion
	Scale    vector.Vector3
}
