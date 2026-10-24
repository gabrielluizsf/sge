package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// drawn is the game holding m drawn onto a w by h canvas.
func drawn(t testing.TB, m *Map, w, h int) *canvas.Canvas {
	g := game.New()
	g.Add(m.Object)
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		t.Fatal(err)
	}
	g.Draw(cv)
	return cv
}

// A map with few tiles draws them one by one rather than its whole
// picture, and the pixels come out the same, wherever it is placed.
func TestASparseMapDrawsAsItsPicture(t *testing.T) {
	for _, at := range []vector.Vector3{{}, {X: 3.4, Y: 1.6}, {X: -2.5, Y: 7.5}} {
		m, _ := New("things", tiles(t), 10, 10)
		m.Set(1, 1, 1)
		m.Set(8, 3, 0)
		m.Set(9, 9, 1)
		m.Transform.Position = at
		if m.Sprite != nil {
			t.Fatal("a map with 3 of 100 cells drawn as a whole picture")
		}
		got := drawn(t, m, 48, 48)
		want, _ := canvas.NewCanvas(48, 48)
		want.Clear(game.New().Background)
		m.Picture()
		m.picture.Draw(want, m.Transform)
		for y := range 48 {
			for x := range 48 {
				if got.At(x, y) != want.At(x, y) {
					t.Fatalf("at %v, pixel %d,%d: %v, not %v", at, x, y, got.At(x, y), want.At(x, y))
				}
			}
		}
	}
}

// The camera moves a sparse map as it moves any picture.
func TestASparseMapFollowsTheCamera(t *testing.T) {
	m, _ := New("things", tiles(t), 10, 10)
	m.Set(2, 2, 1)
	g := game.New()
	g.Add(m.Object)
	g.Camera = vector.Vector3{X: 5, Y: 3}
	cv, _ := canvas.NewCanvas(16, 16)
	g.Draw(cv)
	if cv.At(3, 5) != wall || cv.At(2, 5) == wall || cv.At(3, 4) == wall {
		t.Errorf("the tile at 8,8 seen from 5,3 is not at 3,5: %v %v %v", cv.At(3, 5), cv.At(2, 5), cv.At(3, 4))
	}
}

func TestAFullMapDrawsItsPicture(t *testing.T) {
	m, _ := New("ground", tiles(t), 4, 4)
	m.Fill(0, 0, 4, 2, 0)
	if m.Sprite == nil || m.Sprite.Image != m.Picture() {
		t.Fatal("a half full map is not drawn as its picture")
	}
	m.Set(0, 0, Empty)
	if m.Sprite != nil {
		t.Fatal("a map under half full is still drawn as its picture")
	}
}

// BenchmarkDrawSparseMap draws a 30 by 20 layer of 32 pixel tiles with a
// tile in one cell of ten, as a layer of things over the ground is.
func BenchmarkDrawSparseMap(b *testing.B) {
	sheet, _ := canvas.NewCanvas(64, 32)
	sheet.FillRect(canvas.Area{Width: 32, Height: 32}, grass)
	s, _ := NewTileSet(sheet, 32, 32)
	m, _ := New("things", s, 30, 20)
	for i := 0; i < 600; i += 10 {
		m.Set(i%30, i/30, 0)
	}
	g := game.New()
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(960, 640)
	b.ResetTimer()
	for range b.N {
		g.Draw(cv)
	}
}

// A sparse map has no picture until it needs one: a world made of many
// layers of chunks, mostly empty, would hold a megabyte each otherwise.
func TestASparseMapHoldsNoPicture(t *testing.T) {
	m, _ := New("things", tiles(t), 100, 100)
	m.Set(3, 3, 1)
	if m.picture.Image != nil {
		t.Fatal("a map with one tile allocated its picture")
	}
	m.Fill(0, 0, 100, 50, 0)
	if m.Sprite == nil || m.Picture().At(3*4+1, 3*4+1) != grass || m.Picture().At(1, 60*4) != 0 {
		t.Fatal("once half full, the picture does not show the cells")
	}
	if _, err := New("huge", tiles(t), 1<<16, 1<<16); err == nil {
		t.Error("a map too big for a canvas was made")
	}
}
