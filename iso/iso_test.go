package iso

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// A cell's corners are a diamond's: its far corner the top, its near
// corner the bottom, the others its sides.
func TestACellIsADiamond(t *testing.T) {
	p := New(32, 64, 32)
	for _, c := range []struct{ world, screen vector.Vector3 }{
		{vector.Vector3{}, vector.Vector3{}},
		{vector.Vector3{X: 32}, vector.Vector3{X: 32, Y: 16}},
		{vector.Vector3{Y: 32}, vector.Vector3{X: -32, Y: 16}},
		{vector.Vector3{X: 32, Y: 32}, vector.Vector3{Y: 32}},
		{vector.Vector3{X: 32, Y: 32, Z: 32}, vector.Vector3{Y: 0}},
	} {
		if got := p.Project(c.world); got != c.screen {
			t.Errorf("Project(%v) = %v, want %v", c.world, got, c.screen)
		}
	}
}

// What is drawn at a point of the screen is the ground there.
func TestUnprojectUndoesProject(t *testing.T) {
	p := New(32, 64, 32)
	for _, w := range []vector.Vector3{{X: 5, Y: 7}, {X: -40, Y: 12.5}, {X: 300, Y: -2}} {
		if got := p.Unproject(p.Project(w)); got.Sub(w).Length() > 1e-9 {
			t.Errorf("Unproject(Project(%v)) = %v", w, got)
		}
	}
}
