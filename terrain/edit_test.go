package terrain

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// A tree cut down stays cut: in the chunk shown, and in the chunk made
// again after the camera went away and came back.
func TestAnEditOutlivesItsChunk(t *testing.T) {
	gen := meadow(t)
	g := game.New()
	w := New("world", gen, 8, 0, 1, 2)
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(64, 64)
	g.Draw(cv)
	g.Step(1.0 / 60)
	col, row := -1, 0
	for c := range 16 {
		for r := range 16 {
			if gen.Tile(2, c, r) == tree {
				col, row = c, r
			}
		}
	}
	if col < 0 {
		t.Fatal("no tree in view")
	}
	w.Set(2, col, row, tilemap.Empty)
	ch := w.chunks[mustKey(col/8, row/8)]
	if gen.Tile(2, col, row) != tilemap.Empty || ch.layers[2].At(col%8, row%8) != tilemap.Empty {
		t.Fatal("the tree still stands")
	}
	if !gen.In("trees", col, row) {
		t.Error("the woods moved")
	}
	g.Camera = vector.Vector3{X: 4000}
	g.Step(1.0 / 60)
	if _, ok := w.chunks[mustKey(col/8, row/8)]; ok {
		t.Fatal("the chunk was kept")
	}
	g.Camera = vector.Vector3{}
	g.Step(1.0 / 60)
	if w.chunks[mustKey(col/8, row/8)].layers[2].At(col%8, row%8) != tilemap.Empty {
		t.Error("the tree grew back in the chunk made again")
	}
	for r := -20; r < 20; r++ {
		for c := -20; c < 20; c++ {
			for l := range 3 {
				if w.TileAt(l, c, r) != gen.Tile(l, c, r) {
					t.Fatalf("TileAt %d,%d layer %d is %d, not %d", c, r, l, w.TileAt(l, c, r), gen.Tile(l, c, r))
				}
			}
		}
	}
	if c, r := w.CellOf(-0.5, 4.5); c != -1 || r != 1 {
		t.Errorf("CellOf is %d,%d", c, r)
	}
}

// A chunk's props come into the game with it and leave with it.
func TestPropsComeAndGoWithTheirChunk(t *testing.T) {
	g := game.New()
	w := New("world", meadow(t), 8, 0, 1, 2)
	made := map[[2]int]*game.Object{}
	w.Props = func(cx, cy int) []*game.Object {
		if cx != 0 || cy != 0 {
			return nil
		}
		o := game.NewObject("house")
		made[[2]int{cx, cy}] = o
		return []*game.Object{o}
	}
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(64, 64)
	g.Draw(cv)
	g.Step(1.0 / 60)
	if g.Find("house") == nil {
		t.Fatal("the house was not added")
	}
	g.Camera = vector.Vector3{X: 4000}
	g.Step(1.0 / 60)
	g.Step(1.0 / 60)
	if g.Find("house") != nil {
		t.Error("the house stayed when its chunk went")
	}
}

func mustKey(cx, cy int) int {
	k, _ := cellKey(cx, cy)
	return k
}
