package tilemap_test

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// isoTiles is a set of two 64 by 64 tiles, of cells 32 square: a green
// one and a red one, each a block filling its picture.
func isoTiles(t *testing.T) *tilemap.TileSet {
	sheet, _ := canvas.NewCanvas(128, 64)
	sheet.FillRect(canvas.Area{X: 0, Y: 0, Width: 64, Height: 64}, canvas.RGB(0, 255, 0))
	sheet.FillRect(canvas.Area{X: 64, Y: 0, Width: 64, Height: 64}, canvas.RGB(255, 0, 0))
	ts, err := tilemap.NewTileSet(sheet, 64, 64)
	if err != nil {
		t.Fatal(err)
	}
	ts.CellWidth, ts.CellHeight = 32, 32
	return ts
}

// A projected map stands each tile's picture on the bottom corner of its
// cell's diamond, and keeps the cells of the world square.
func TestIsoMapStandsTilesOnTheirCells(t *testing.T) {
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	m, _ := tilemap.New("ground", isoTiles(t), 2, 2)
	m.Set(1, 1, 0)
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(200, 200)
	g.Camera = g.ToWorld(vector.Vector3{X: -100, Y: -20}) // the map's origin 100, 20 in
	g.Draw(cv)
	bottom := g.ToScreen(vector.Vector3{X: 64, Y: 64}) // cell 1, 1's near corner
	if got := cv.At(int(bottom.X), int(bottom.Y)-1); got != canvas.RGB(0, 255, 0) {
		t.Errorf("no tile above the cell's near corner: %v", got)
	}
	if got := cv.At(int(bottom.X), int(bottom.Y)+1); got == canvas.RGB(0, 255, 0) {
		t.Errorf("the tile hangs below its cell")
	}
	if c, r, ok := m.CellAt(40, 40); !ok || c != 1 || r != 1 {
		t.Errorf("CellAt(40, 40) = %d, %d, %v: cells are 32 square in the world", c, r, ok)
	}
}

// A Sorted map's tiles stand among the objects of their Z: an object in
// front of a tile's cell is drawn over it, one behind it under it.
func TestSortedMapTilesStandAmongObjects(t *testing.T) {
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	m, _ := tilemap.New("trees", isoTiles(t), 3, 3)
	m.Sorted = true
	m.Set(1, 1, 1) // red, at the middle cell
	g.Add(m.Object)
	img, _ := canvas.NewCanvas(4, 4)
	img.Clear(canvas.RGB(0, 0, 255))
	near := game.NewObject("near")
	near.Sprite = &sprite.Sprite{Image: img, PivotX: 0.5, PivotY: 1}
	near.Transform.Position = vector.Vector3{X: 70, Y: 70} // in front of cell 1, 1
	g.Add(near)
	cv, _ := canvas.NewCanvas(300, 300)
	g.Camera = g.ToWorld(vector.Vector3{X: -150, Y: -60})
	g.Draw(cv)
	at := g.ToScreen(near.Transform.Position)
	if got := cv.At(int(at.X), int(at.Y)-2); got != canvas.RGB(0, 0, 255) {
		t.Errorf("the object in front of the tile is behind it: %v", got)
	}
	near.Transform.Position = vector.Vector3{X: 34, Y: 34} // behind it
	g.Draw(cv)
	at = g.ToScreen(near.Transform.Position)
	if got := cv.At(int(at.X), int(at.Y)-2); got != canvas.RGB(255, 0, 0) {
		t.Errorf("the object behind the tile is in front of it: %v", got)
	}
}
