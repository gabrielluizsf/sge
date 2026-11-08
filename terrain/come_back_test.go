package terrain

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// A place seen, left far behind until its chunks are let go, and come
// back to, is drawn exactly as it was, and the chunks kept stay few.
func TestAPlaceComesBackIdentical(t *testing.T) {
	g := game.New()
	g.Covered = true
	w := New("world", meadow(t), 8, 0, 1, 2)
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(96, 64)
	g.Draw(cv)
	look := func(at vector.Vector3) []canvas.Color {
		g.Camera = at
		for range 3 {
			g.Step(1.0 / 60)
			g.Draw(cv)
		}
		return slices.Clone(cv.Pixels)
	}
	home := vector.Vector3{X: 40, Y: -30}
	before := look(home)
	for i := 1; i <= 40; i++ { // walk a long way off, chunk after chunk
		look(vector.Vector3{X: 40 + float64(i)*200, Y: -30 + float64(i)*90})
		if n := len(w.chunks); n > 4*4*2 {
			t.Fatalf("%d chunks kept at step %d", n, i)
		}
	}
	if _, kept := w.chunks[mustKey(1, -1)]; kept {
		t.Fatal("the chunks at home were kept")
	}
	if !slices.Equal(before, look(home)) {
		t.Fatal("home looks different after coming back")
	}
}
