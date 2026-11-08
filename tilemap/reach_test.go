package tilemap_test

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// A tile whose picture is larger than the tile set's is placed by its own
// bottom middle, standing on its cell, and is drawn while its crown is in
// view though its cell is out of it, to the side as well as below.
func TestWideTilesReachIntoTheViewFromOutsideIt(t *testing.T) {
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	ts := tilemap.NewLazyTileSet(tilemap.Lazy{Count: 1, Width: 64, Height: 64, Make: func(int) *canvas.Canvas {
		cv, _ := canvas.NewCanvas(240, 300)
		cv.Clear(canvas.RGB(0, 255, 0))
		return cv
	}})
	ts.CellWidth, ts.CellHeight = 32, 32
	ts.Reach = 300
	m, _ := tilemap.New("trees", ts, 40, 40)
	m.Set(20, 20, 0)
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(100, 100)
	bottom := g.ToScreen(vector.Vector3{X: 21 * 32, Y: 21 * 32}) // the cell's near corner
	// The view's left edge is 60 pixels right of the tree's foot: its crown, 120 each way, reaches in.
	g.Camera = g.ToWorld(vector.Vector3{X: bottom.X + 90, Y: bottom.Y - 60})
	g.Draw(cv)
	green := canvas.RGB(0, 255, 0)
	if got := cv.At(5, 20); got != green {
		t.Errorf("the crown does not reach into the view from the side: %v", got)
	}
	if got := cv.At(5, 99); got == green {
		t.Errorf("the tile hangs below its cell: %v", got)
	}
}
