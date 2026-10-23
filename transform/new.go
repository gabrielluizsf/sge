package transform

import "github.com/gabrielluizsf/sge/vector"

// New is a transform at the origin, not rotated, at scale 1.
func New() Transform {
	return Transform{
		Rotation: vector.Identity(),
		Scale:    vector.Vector3{X: 1, Y: 1, Z: 1},
	}
}
