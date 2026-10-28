package game

import (
	"testing"

	"github.com/gabrielluizsf/sge/collider"
)

// counting is a shape that counts the times it is asked whether it
// overlaps another, and always does.
type counting struct{ asked *int }

func (c counting) Bounds() collider.Box                   { return collider.Box{Width: 10, Height: 10} }
func (c counting) Moved(dx, dy float64) collider.Collider { return c }
func (c counting) Overlaps(other collider.Collider) bool  { *c.asked++; return true }

// Two objects neither of which handles collisions are not tested against
// each other: no one would be told. Walls and the ground are many such.
func TestObjectsNoOneHandlesAreNotTested(t *testing.T) {
	asked := 0
	g := New()
	for _, name := range []string{"wall", "floor"} {
		g.Add(NewObject(name)).Collider = counting{&asked}
	}
	g.Step(0.1)
	if asked != 0 {
		t.Errorf("tested %d times", asked)
	}
	var log []string
	hero := g.Add(NewObject("hero"))
	hero.Collider = collider.Box{Width: 1, Height: 1}
	hero.AddScript(recorder{&log})
	g.Step(0.1)
	if asked == 0 || len(log) < 3 {
		t.Errorf("the hero, which handles collisions, was not tested: asked %d, log %v", asked, log)
	}
}
