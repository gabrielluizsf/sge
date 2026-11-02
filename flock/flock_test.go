package flock

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// trees are perches on a line, 200 apart.
func trees(vector.Vector3, float64) []Perch {
	var out []Perch
	for i := range 8 {
		out = append(out, Perch{At: vector.Vector3{X: float64(i) * 200}, Height: 100})
	}
	return out
}

func newFlock(scare *vector.Vector3) *Flock {
	f := New(7, 2)
	f.Perches = trees
	f.Around = func() vector.Vector3 { return vector.Vector3{X: 700} }
	f.Scares = func() []vector.Vector3 {
		if scare == nil {
			return nil
		}
		return []vector.Vector3{*scare}
	}
	return f
}

// Birds sit on perches until someone comes near, then fly and land on
// another.
func TestBirdsRiseWhenSomeoneComesNearAndLandElsewhere(t *testing.T) {
	var who *vector.Vector3
	f := newFlock(who)
	f.Step(0.1)
	b := f.Birds()[0]
	if b.State != Perched || b.Alt != 100 {
		t.Fatalf("bird starts %v at %v", b.State, b.Alt)
	}
	start := b.Pos
	near := start.Add(vector.Vector3{X: 30})
	who = &near
	f = newFlock(who)
	f.Step(0.1)
	f.Step(0.1)
	b = f.Birds()[0]
	if b.State != Flying {
		t.Fatalf("bird did not rise: %v", b.State)
	}
	away := vector.Vector3{X: 1e6}
	who = &away
	f.Scares = func() []vector.Vector3 { return []vector.Vector3{away} }
	flew := false
	for range 400 {
		f.Step(0.05)
		flew = flew || b.Alt > 100
		if b.State == Perched {
			break
		}
	}
	if !flew || b.State != Perched || b.Pos == start {
		t.Fatalf("bird flew %v, ended %v at %v from %v", flew, b.State, b.Pos, start)
	}
}

// The same seed is the same flock.
func TestTheSameSeedMakesTheSameFlock(t *testing.T) {
	run := func() vector.Vector3 {
		f := newFlock(nil)
		for range 2000 {
			f.Step(0.05)
		}
		return f.Birds()[1].Pos
	}
	if run() != run() {
		t.Fatal("two flocks of one seed differ")
	}
}
