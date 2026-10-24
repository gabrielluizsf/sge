package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// A map out of view does not animate, which no one would see; it shows
// the right frame again on the step it comes into view.
func TestAMapOutOfViewIsNotAnimatedUntilSeen(t *testing.T) {
	m, _ := New("field", swaying(t, false), 2, 1)
	m.Fill(0, 0, 2, 1, 0)
	g := game.New()
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(8, 4)
	g.Draw(cv)
	g.Camera = vector.Vector3{X: 100}
	g.Step(0.3) // the second frame, were it seen
	if m.painted[0] != 1 || m.Picture().At(1, 1) != still {
		t.Fatalf("animated out of view: showing %d", m.painted[0])
	}
	g.Camera = vector.Vector3{}
	g.Step(0.01)
	if m.Picture().At(1, 1) != blade {
		t.Errorf("back in view, showing %v, not the frame of the time", m.Picture().At(1, 1))
	}
}

// A tile is copied into the map's picture as it is, its transparency too.
func TestPaintKeepsATilesTransparency(t *testing.T) {
	sheet, _ := canvas.NewCanvas(4, 4)
	half := canvas.RGBA(canvas.RGBAParams{R: 200, A: 128})
	sheet.FillRect(canvas.Area{Width: 4, Height: 4}, half)
	s, _ := NewTileSet(sheet, 4, 4)
	m, _ := New("glass", s, 1, 1)
	m.Set(0, 0, 0)
	if got := m.Picture().At(1, 1); got.A() != sheet.At(1, 1).A() || got != sheet.At(1, 1) {
		t.Errorf("painted %v, the tile is %v", got, sheet.At(1, 1))
	}
}
