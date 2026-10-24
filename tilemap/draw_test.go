package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
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

// However full a map is, of opaque tiles or not, it draws the pixels its
// whole picture would, wherever it is placed: tile by tile, or copied
// whole when every cell is opaque.
func TestAMapDrawsAsItsPicture(t *testing.T) {
	fills := map[string]func(m *Map){
		"sparse": func(m *Map) { m.Set(1, 1, 1); m.Set(8, 3, 0); m.Set(9, 9, 1) },
		"half":   func(m *Map) { m.Fill(0, 0, 10, 6, 0); m.Set(2, 2, Empty) },
		"opaque": func(m *Map) { m.Fill(0, 0, 10, 10, 0); m.Set(4, 4, 1) },
		"glass":  func(m *Map) { m.Fill(0, 0, 10, 10, 2) },
		"mixed":  func(m *Map) { m.Fill(0, 0, 10, 4, 2); m.Fill(0, 4, 10, 3, 1) },
	}
	for name, fill := range fills {
		for _, at := range []vector.Vector3{{}, {X: 3.4, Y: 1.6}, {X: -2.5, Y: 7.5}} {
			m, _ := New("things", glassy(t), 10, 10)
			fill(m)
			m.Transform.Position = at
			got := drawn(t, m, 48, 48)
			want, _ := canvas.NewCanvas(48, 48)
			want.Clear(game.New().Background)
			(&sprite.Sprite{Image: m.Picture()}).Draw(want, m.Transform)
			for y := range 48 {
				for x := range 48 {
					if got.At(x, y) != want.At(x, y) {
						t.Fatalf("%s at %v, pixel %d,%d: %v, not %v", name, at, x, y, got.At(x, y), want.At(x, y))
					}
				}
			}
		}
	}
}

// glassy is tiles with a tile 2 half transparent.
func glassy(t *testing.T) *TileSet {
	sheet, _ := canvas.NewCanvas(12, 4)
	sheet.FillRect(canvas.Area{Width: 4, Height: 4}, grass)
	sheet.FillRect(canvas.Area{X: 4, Width: 4, Height: 4}, wall)
	// Put, as filling would blend onto the transparent sheet.
	for y := range 4 {
		for x := 8; x < 12; x++ {
			sheet.Put(x, y, canvas.RGBA(canvas.RGBAParams{R: 200, A: 100}))
		}
	}
	s, err := NewTileSet(sheet, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	return s
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

// Only a map that is copied whole needs its picture: one wholly of opaque
// tiles. One with a cell empty or see-through draws tile by tile.
func TestOnlyAnOpaqueMapIsCopiedWhole(t *testing.T) {
	m, _ := New("ground", tiles(t), 4, 4)
	m.Fill(0, 0, 4, 4, 0)
	if !m.opaque() {
		t.Fatal("a map of grass and walls is not opaque")
	}
	m.Set(0, 0, Empty)
	if m.opaque() {
		t.Fatal("a map with an empty cell is opaque")
	}
	g, _ := New("glass", glassy(t), 2, 1)
	g.Fill(0, 0, 2, 1, 2)
	if g.opaque() {
		t.Fatal("a map of see-through tiles is opaque")
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

// A map has no picture until it needs one: a world made of many layers
// of chunks, mostly empty, would hold a megabyte each otherwise.
func TestAMapHoldsNoPictureUntilCopiedWhole(t *testing.T) {
	m, _ := New("things", tiles(t), 100, 100)
	m.Set(3, 3, 1)
	m.Fill(0, 0, 100, 50, 0)
	drawn(t, m, 8, 8)
	if m.picture.Image != nil {
		t.Fatal("a map drawn tile by tile allocated its picture")
	}
	m.Fill(0, 50, 100, 50, 1)
	drawn(t, m, 8, 8)
	if m.picture.Image == nil || m.Picture().At(3*4+1, 3*4+1) != grass || m.Picture().At(1, 60*4+1) != wall {
		t.Fatal("copied whole, the picture does not show the cells")
	}
	if _, err := New("huge", tiles(t), 1<<16, 1<<16); err == nil {
		t.Error("a map too big for a canvas was made")
	}
}
