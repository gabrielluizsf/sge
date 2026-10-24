package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

var (
	blade = canvas.RGB(60, 200, 80)
	still = canvas.RGB(30, 120, 40)
)

// swaying is a set whose tile 0 sways: it shows tiles 1 and 2 in turn,
// four times a second. Tile 3 is plain.
func swaying(t *testing.T, staggered bool) *TileSet {
	sheet, _ := canvas.NewCanvas(16, 4)
	sheet.FillRect(canvas.Area{X: 4, Width: 4, Height: 4}, still)
	sheet.FillRect(canvas.Area{X: 8, Width: 4, Height: 4}, blade)
	sheet.FillRect(canvas.Area{X: 12, Width: 4, Height: 4}, wall)
	s, err := NewTileSet(sheet, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	s.Animate(0, Animation{Frames: []int{1, 2}, FPS: 4, Staggered: staggered})
	return s
}

func TestAnimatedTilesTurnOverAsTheWorldSteps(t *testing.T) {
	m, _ := New("field", swaying(t, false), 3, 1)
	m.Fill(0, 0, 2, 1, 0)
	m.Set(2, 0, 3)
	g := game.New()
	g.Add(m.Object)
	img := m.Picture()
	if img.At(1, 1) != still || img.At(5, 1) != still {
		t.Fatalf("first frame: %v %v", img.At(1, 1), img.At(5, 1))
	}
	g.Step(0.3) // a quarter of a second in: the second frame
	if img.At(1, 1) != blade || img.At(5, 1) != blade || img.At(9, 1) != wall {
		t.Errorf("after 0.3s: %v %v %v", img.At(1, 1), img.At(5, 1), img.At(9, 1))
	}
	g.Step(0.25)
	if img.At(1, 1) != still {
		t.Errorf("after 0.55s: %v", img.At(1, 1))
	}
	if m.At(0, 0) != 0 {
		t.Errorf("the cell holds tile %d, not the animated tile", m.At(0, 0))
	}
	m.Set(0, 0, 3)
	g.Step(0.25)
	if img.At(1, 1) != wall {
		t.Error("a cell set to a plain tile kept animating")
	}
}

func TestStaggeredTilesDoNotAllTurnTogether(t *testing.T) {
	m, _ := New("meadow", swaying(t, true), 8, 8)
	m.Fill(0, 0, 8, 8, 0)
	frames := map[canvas.Color]int{}
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			frames[m.Picture().At(c*4, r*4)]++
		}
	}
	if frames[still] == 0 || frames[blade] == 0 {
		t.Errorf("all 64 cells show one frame: %v", frames)
	}
}
