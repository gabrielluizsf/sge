package flock

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/vector"
)

// State is what a bird is doing.
type State int

// The states: sitting on a perch, or flying to another.
const (
	Perched State = iota
	Flying
)

// Bird is one bird: where it is on the ground, how high above it in
// pixels, what it is doing, the way it faces (a unit vector along the
// ground), and Time, the seconds it has been in its state, to pick a
// frame of its wings by.
type Bird struct {
	Pos   vector.Vector3
	Alt   float64
	State State
	Dir   vector.Vector3
	Time  float64

	from, to Perch
	leg      float64 // how far along the flight, 0 to 1
	span     float64 // how long the flight is, in world units
	rest     float64 // seconds left perched before it hops off by itself
	src      *random.Source
}
