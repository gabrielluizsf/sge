package tilemap

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

var (
	grass = canvas.RGB(40, 160, 60)
	wall  = canvas.RGB(120, 120, 120)
)

// tiles is a set of two 4x4 tiles, grass and a solid wall.
func tiles(t *testing.T) *TileSet {
	sheet, err := canvas.NewCanvas(8, 4)
	if err != nil {
		t.Fatal(err)
	}
	sheet.FillRect(canvas.Area{Width: 4, Height: 4}, grass)
	sheet.FillRect(canvas.Area{X: 4, Width: 4, Height: 4}, wall)
	s, err := NewTileSet(sheet, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	s.SetSolid(true, 1)
	return s
}

const level = `
#####
#...#
#.#.#
#####`

func parse(t *testing.T) *Map {
	m, err := Parse("level", tiles(t), level, map[rune]int{'.': 0, '#': 1})
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestTileSetCutsTilesAndKnowsTheSolidOnes(t *testing.T) {
	s := tiles(t)
	if len(s.Tiles) != 2 || s.Tiles[1].Image.At(0, 0) != wall {
		t.Fatalf("tiles = %d", len(s.Tiles))
	}
	if s.IsSolid(0) || !s.IsSolid(1) || s.IsSolid(Empty) || s.IsSolid(7) {
		t.Error("solid tiles wrong")
	}
}

func TestParseLaysOutCellsAndPaintsThem(t *testing.T) {
	m := parse(t)
	if m.Cols != 5 || m.Rows != 4 || m.At(1, 1) != 0 || m.At(2, 2) != 1 || m.At(9, 9) != Empty {
		t.Fatalf("map %dx%d", m.Cols, m.Rows)
	}
	img := m.Picture()
	if img.Width != 20 || img.Height != 16 || img.At(5, 5) != grass || img.At(9, 9) != wall {
		t.Errorf("picture %dx%d: %v %v", img.Width, img.Height, img.At(5, 5), img.At(9, 9))
	}
	m.Set(2, 2, Empty)
	if m.At(2, 2) != Empty || img.At(9, 9) != 0 {
		t.Errorf("cleared cell still painted: %v", img.At(9, 9))
	}
	m.Fill(1, 1, 3, 2, 1)
	if m.At(3, 2) != 1 || img.At(13, 9) != wall {
		t.Error("Fill did not fill")
	}
}

func TestCellsAndWorldPointsFollowTheMapsPositionAndScale(t *testing.T) {
	m := parse(t)
	m.Transform.Position = vector.Vector3{X: 100, Y: 50}
	m.Transform.Scale = vector.Vector3{X: 2, Y: 2, Z: 1}
	if c, r, ok := m.CellAt(100+8*2+1, 50+8*1+1); !ok || c != 2 || r != 1 {
		t.Errorf("CellAt = %d, %d, %v", c, r, ok)
	}
	if _, _, ok := m.CellAt(99, 50); ok {
		t.Error("a point left of the map is on it")
	}
	if got := m.CellCenter(2, 1); got.X != 100+16+4 || got.Y != 50+8+4 {
		t.Errorf("CellCenter = %v", got)
	}
	if !m.SolidAt(101, 51) || m.SolidAt(100+8+1, 50+8+1) {
		t.Error("SolidAt wrong")
	}
	if m.TileAt(101, 51) != 1 || m.TileAt(100+8+1, 50+8+1) != 0 || m.TileAt(99, 50) != Empty {
		t.Error("TileAt wrong")
	}
}

func TestObjectsCollideWithSolidTilesOnly(t *testing.T) {
	m := parse(t)
	g := game.New()
	g.Add(m.Object)
	hits := map[string]bool{}
	for _, tc := range []struct {
		name string
		at   vector.Vector3
	}{{"on grass", vector.Vector3{X: 6, Y: 6}}, {"on the middle wall", vector.Vector3{X: 10, Y: 10}}, {"at the edge", vector.Vector3{X: 2, Y: 6}}} {
		o := g.Add(game.NewObject(tc.name))
		o.Collider = collider.Circle{Radius: 1}
		o.Transform.Position = tc.at
		o.AddScript(recorder{hits})
	}
	g.Step(0)
	if hits["on grass"] || !hits["on the middle wall"] || !hits["at the edge"] {
		t.Errorf("hits = %v", hits)
	}
}

type recorder struct{ hits map[string]bool }

func (r recorder) Update(o *game.Object, dt float64) {}
func (r recorder) Collision(o, other *game.Object) {
	if other.Name == "level" {
		r.hits[o.Name] = true
	}
}

func TestMoveSlidesAlongWalls(t *testing.T) {
	m := parse(t)
	o := game.NewObject("walker")
	o.Collider = collider.Box{X: -1, Y: -1, Width: 2, Height: 2}
	o.Transform.Position = vector.Vector3{X: 6, Y: 6}
	hitX, hitY := m.Move(o, -10, 1)
	p := o.Transform.Position
	if !hitX || hitY || p.X < 5 || p.X > 5.02 || p.Y != 7 {
		t.Errorf("Move = %v %v, at %v: want stopped at the left wall's edge, x 5", hitX, hitY, p)
	}
	free := game.NewObject("ghost")
	m.Move(free, -100, 0)
	if free.Transform.Position.X != -100 {
		t.Error("an object without a collider was stopped")
	}
}

func BenchmarkDrawMap(b *testing.B) {
	sheet, _ := canvas.NewCanvas(32, 16)
	set, _ := NewTileSet(sheet, 16, 16)
	m, _ := New("big", set, 30, 20)
	m.Fill(0, 0, 30, 20, 0)
	g := game.New()
	g.Add(m.Object)
	cv, _ := canvas.NewCanvas(480, 320)
	for b.Loop() {
		g.Draw(cv)
	}
}
