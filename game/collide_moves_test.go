package game

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/vector"
)

// mover is a collision handler that changes b when it hits it: moves it
// away, or swaps its collider.
type mover struct {
	log    *[]string
	change func(other *Object)
}

func (m mover) Update(o *Object, dt float64) {}
func (m mover) Collision(o, other *Object) {
	*m.log = append(*m.log, o.Name+" hit "+other.Name)
	if other.Name == "b" {
		m.change(other)
	}
}

// Objects a collision handler moves, or gives another shape, collide
// where and as they are now for the rest of the step.
func TestCollisionsSeeChangesMadeDuringTheStep(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(other *Object)
	}{
		{"moved", func(b *Object) { b.Transform.Position = vector.Vector3{X: 1000} }},
		{"reshaped", func(b *Object) { b.Collider = collider.Box{X: 500, Y: 500, Width: 1, Height: 1} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var log []string
			g := New()
			a := g.Add(NewObject("a"))
			a.Collider = collider.Circle{Radius: 5}
			a.AddScript(mover{&log, tc.change})
			b := g.Add(NewObject("b"))
			b.Collider = collider.Circle{Radius: 5}
			b.AddScript(recorder{&log})
			c := g.Add(NewObject("c"))
			c.Collider = collider.Circle{Radius: 5}
			c.AddScript(recorder{&log})
			g.Step(0)
			if slices.Contains(log, "b hit c") || slices.Contains(log, "c hit b") {
				t.Errorf("log = %v: b was changed away from c by a, yet collided with it", log)
			}
			if !slices.Contains(log, "a hit c") {
				t.Errorf("log = %v, want a to hit c", log)
			}
		})
	}
}
