package tilemap_test

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// revealWorld is a 4 by 4 projected map of a tall red block at cell 2, 2,
// drawn onto a canvas; keep, when set, is what the map keeps in sight.
func revealWorld(t *testing.T, keep []game.Reveal) (*canvas.Canvas, vector.Vector3) {
	sheet, _ := canvas.NewCanvas(64, 96)
	sheet.FillRect(canvas.Area{Width: 64, Height: 96}, canvas.RGB(255, 0, 0))
	ts, err := tilemap.NewTileSet(sheet, 64, 96)
	if err != nil {
		t.Fatal(err)
	}
	ts.Height, ts.CellWidth, ts.CellHeight = 32, 32, 32
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	m, _ := tilemap.New("trees", ts, 4, 4)
	m.Sorted = true
	m.Set(2, 2, 0)
	if keep != nil {
		g.Reveals = func() []game.Reveal { return keep }
	}
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(300, 300)
	g.Camera = g.ToWorld(vector.Vector3{X: -150, Y: -100})
	g.Draw(cv)
	return cv, g.ToScreen(m.CellCenter(1, 1))
}

// A tall tile that covers what is kept in sight behind it is drawn faded,
// and stays whole when nothing is behind it.
func TestATallTileFadesOverWhatIsKeptInSight(t *testing.T) {
	cv, at := revealWorld(t, nil)
	x, y := int(at.X), int(at.Y)-10
	if got := cv.At(x, y); got != canvas.RGB(255, 0, 0) {
		t.Fatalf("whole tile is %v", got)
	}
	behind := vector.Vector3{X: 48, Y: 48}
	cv, at = revealWorld(t, []game.Reveal{{At: behind, Height: 40}})
	got := cv.At(int(at.X), int(at.Y)-10)
	if got.R() > 160 || got.R() == 0 {
		t.Errorf("tile over the farmer is %v, want faded", got)
	}
	// A farmer in front of the tile does not fade it.
	front := vector.Vector3{X: 112, Y: 112}
	cv, at = revealWorld(t, []game.Reveal{{At: front, Height: 40}})
	if cv.At(int(at.X), int(at.Y)-10) != canvas.RGB(255, 0, 0) {
		t.Errorf("tile fades over a farmer in front of it")
	}
}
