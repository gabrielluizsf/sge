package flock

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/vector"
)

// Flock is a few birds kept about a view. Count of them live near the
// point Around gives; those that fall further than Range from it are put
// on a perch near it again. Perches says where birds may sit within a
// radius of a point; Scares where those are that make them rise, a
// player, villagers, when one comes within Fear. Speed is how fast they
// fly, in world units a second, Lift how high, in pixels, they climb in
// flight above a straight line between perches.
type Flock struct {
	Seed        uint32
	Count       int
	Range, Fear float64
	Speed, Lift float64
	Perches     func(c vector.Vector3, radius float64) []Perch
	Scares      func() []vector.Vector3
	Around      func() vector.Vector3
	birds       []*Bird
	seated      int // how many birds have been put on perches, for their seeds
}

// New is a flock of count birds from seed, ready to be given its
// functions.
func New(seed uint32, count int) *Flock {
	return &Flock{Seed: seed, Count: count, Range: 640, Fear: 110, Speed: 220, Lift: 90}
}

// Birds are the birds now.
func (f *Flock) Birds() []*Bird { return f.birds }

// source is a new bird's own random source.
func (f *Flock) source() *random.Source {
	f.seated++
	return random.New(f.Seed*2654435761 + uint32(f.seated))
}
