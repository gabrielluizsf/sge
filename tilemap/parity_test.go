package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// A tile of parity shows each cell the frame its place among the four
// asks for, and never turns over.
func TestParityTilesShowThePatchOfTheirPlace(t *testing.T) {
	sheet, _ := canvas.NewCanvas(20, 4)
	var colours [5]canvas.Color
	for i := range colours {
		colours[i] = canvas.RGB(uint8(40*i+20), 90, 200)
		sheet.FillRect(canvas.Area{X: 4 * i, Width: 4, Height: 4}, colours[i])
	}
	s, err := NewTileSet(sheet, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	s.Animate(0, Animation{Frames: []int{1, 2, 3, 4}, FPS: 4, Parity: true})
	m, _ := New("meadow", s, 4, 4)
	m.Fill(0, 0, 4, 4, 0)
	g := game.New()
	g.Add(m.Object)
	g.Step(0.6)
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if got, want := m.Picture().At(c*4, r*4), colours[1+c%2+2*(r%2)]; got != want {
				t.Errorf("cell %d,%d shows %v, not %v", c, r, got, want)
			}
		}
	}
}
