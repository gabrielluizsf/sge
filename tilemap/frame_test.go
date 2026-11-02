package tilemap

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/vector"
)

// swayWorld is a projected world of one big map of tiles that sway, or are
// patterned by where they lie, of squares 8 pixels across on cells of 4.
func swayWorld(t *testing.T, parity bool) (*game.Game, *Map) {
	sheet, _ := canvas.NewCanvas(8*4, 8)
	for i, c := range []canvas.Color{canvas.RGB(20, 120, 30), canvas.RGB(40, 160, 50), canvas.RGB(60, 190, 70), canvas.RGB(90, 220, 90)} {
		sheet.FillRect(canvas.Area{X: i * 8, Width: 8, Height: 8}, c)
	}
	ts, err := NewTileSet(sheet, 8, 8)
	if err != nil {
		t.Fatal(err)
	}
	ts.CellWidth, ts.CellHeight = 4, 4
	ts.Animate(0, Animation{Frames: []int{0, 1, 2, 3}, FPS: 3, Staggered: true, Parity: parity})
	m, _ := New("field", ts, 60, 60)
	m.Fill(0, 0, 60, 60, 0)
	g := game.New()
	g.Projection = iso.New(4, 8, 4)
	g.Covered = true
	g.Add(m.Object)
	return g, m
}

// sweep moves the camera by a step across the field, back and forth.
func sweep(g *game.Game, i int) {
	g.Camera = vector.Vector3{X: 60 + float64(i%40)*3, Y: 40 + float64(i%25)*2}
}

// Animating only the cells in view shows what animating them all does: a
// cell that comes into view shows the frame of the time it is drawn at.
func TestAnimatingOnlyTheViewIsAnimatingAll(t *testing.T) {
	a, ma := swayWorld(t, false)
	b, mb := swayWorld(t, false)
	ca, _ := canvas.NewCanvas(48, 32)
	cb, _ := canvas.NewCanvas(48, 32)
	for i := range 160 {
		a.Step(1.0 / 16)
		b.Step(1.0 / 16)
		mb.refresh(mb.whole()) // every cell, as a step used to
		sweep(a, i)            // the camera moves on after the step, as a follower's does
		sweep(b, i)
		a.Draw(ca)
		b.Draw(cb)
		if !slices.Equal(ca.Pixels, cb.Pixels) {
			t.Fatalf("frame %d: the view shows other tiles than animating every cell does", i)
		}
	}
	apart := 0
	for i := range ma.painted {
		if ma.painted[i] != mb.painted[i] {
			apart++
		}
	}
	if apart == 0 {
		t.Error("every cell was animated: nothing was saved")
	}
}

// A swaying field is not still: drawn frame by frame, it changes what is
// on the canvas as time goes, whatever the player does.
func TestASwayingFieldNeverStandsStill(t *testing.T) {
	g, _ := swayWorld(t, false)
	g.Camera = vector.Vector3{X: 100, Y: 80}
	cv, _ := canvas.NewCanvas(48, 32)
	scratch, _ := canvas.NewCanvas(48, 32)
	drawn := 0
	var seen [][]canvas.Color
	for i := range 120 { // two seconds at 60 a second
		if g.Frame(cv, 1.0/60) {
			drawn++
		}
		g.Draw(scratch) // a whole draw of the same world
		if !slices.Equal(cv.Pixels, scratch.Pixels) {
			t.Fatalf("frame %d is not what drawing it whole makes", i)
		}
		if i%30 == 0 {
			seen = append(seen, slices.Clone(cv.Pixels))
		}
	}
	for i := 1; i < len(seen); i++ {
		if slices.Equal(seen[i], seen[i-1]) {
			t.Errorf("the field looks the same half a second after %d", i-1)
		}
	}
	if drawn < 40 {
		t.Errorf("only %d of 120 frames of a swaying field were drawn", drawn)
	}
}

// Tiles patterned by where they lie never turn over: a still field of
// them is not repainted by the clock, nor drawn again frame after frame.
func TestPatternedTilesAreLeftAlone(t *testing.T) {
	g, m := swayWorld(t, true)
	g.Camera = vector.Vector3{X: 100, Y: 80}
	cv, _ := canvas.NewCanvas(48, 32)
	g.Frame(cv, 1.0/60)
	if m.Tiles.turning != 0 {
		t.Fatalf("%d animations turn over", m.Tiles.turning)
	}
	drawn := 0
	for range 120 {
		if g.Frame(cv, 1.0/60) {
			drawn++
		}
	}
	if drawn != 0 {
		t.Errorf("a field of patterned tiles was drawn %d times standing still", drawn)
	}
	if got := m.painted[3*m.Cols+2]; got != 2 { // column 2, row 3: the third frame
		t.Errorf("a patterned tile was not set by where it lies: %d", got)
	}
}
