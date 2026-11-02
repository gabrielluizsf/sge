package flock

import "github.com/gabrielluizsf/sge/vector"

// Perch is a place a bird may sit: a point of the ground, and how high
// above it, in pixels of the screen, the bird sits there.
type Perch struct {
	At     vector.Vector3
	Height float64
}
