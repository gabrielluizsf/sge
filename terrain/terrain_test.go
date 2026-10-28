package terrain

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

var (
	grassInk = canvas.RGB(40, 160, 60)
	waterInk = canvas.RGB(40, 90, 200)
	treeInk  = canvas.RGB(20, 80, 30)
)

const grass, water, tree = 0, 1, 2

// tiles is a set of 4 by 4 tiles: grass, water, which is solid, and a
// tree, a dot on a transparent tile.
func tiles(t testing.TB) *tilemap.TileSet {
	sheet, _ := canvas.NewCanvas(12, 4)
	sheet.FillRect(canvas.Area{Width: 4, Height: 4}, grassInk)
	sheet.FillRect(canvas.Area{X: 4, Width: 4, Height: 4}, waterInk)
	sheet.FillRect(canvas.Area{X: 9, Y: 1, Width: 2, Height: 2}, treeInk)
	s, err := tilemap.NewTileSet(sheet, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	s.SetSolid(true, water)
	return s
}

// meadow is grass everywhere, lakes over it, and trees kept a cell away
// from the water.
func meadow(t testing.TB) *Generator {
	g := &Generator{Tiles: tiles(t), Layers: 3}
	lakes := Blob(7, 8, 0.6)
	trees := Scatter(9, 0.2, 1)
	g.Add(Area{Name: "grass", Layer: 0, In: Everywhere, Tile: func(int, int) int { return grass }})
	g.Add(Area{Name: "water", Layer: 1, In: lakes, Tile: func(int, int) int { return water }})
	g.Add(Area{Name: "trees", Layer: 2, In: func(c, r int) bool { return trees(c, r) && !g.Near("water", c, r, 1) },
		Tile: func(int, int) int { return tree }})
	return g
}

func TestNoiseIsTheSameEveryTimeAndSmooth(t *testing.T) {
	for i := range 200 {
		x, y := float64(i)*0.37-30, float64(i)*0.71-50
		a := Fractal(3, x, y, 3)
		if a != Fractal(3, x, y, 3) || a < 0 || a > 1 {
			t.Fatalf("Fractal(%v, %v) = %v", x, y, a)
		}
		if d := math.Abs(Noise(3, x, y) - Noise(3, x+0.01, y)); d > 0.05 {
			t.Fatalf("noise jumps %v in a hundredth", d)
		}
	}
	if Hash(1, 5, 5) == Hash(2, 5, 5) || Hash(1, 5, 5) == Hash(1, 5, 6) || Hash(1, -5, 5) == Hash(1, 5, 5) {
		t.Error("hashes of other seeds or cells agree")
	}
}

func TestScatterKeepsItsSpacing(t *testing.T) {
	in := Scatter(4, 0.5, 2)
	var got [][2]int
	for r := -40; r < 40; r++ {
		for c := -40; c < 40; c++ {
			if in(c, r) {
				got = append(got, [2]int{c, r})
			}
		}
	}
	for i, a := range got {
		for _, b := range got[i+1:] {
			if max(abs(a[0]-b[0]), abs(a[1]-b[1])) <= 2 {
				t.Fatalf("%v and %v are within 2 cells", a, b)
			}
		}
	}
	if len(got) < 6400/25/3 {
		t.Errorf("only %d of 6400 cells", len(got))
	}
}

func abs(x int) int { return max(x, -x) }

func TestThereAre47Masks(t *testing.T) {
	if n := len(Masks()); n != 47 {
		t.Errorf("%d masks", n)
	}
	in := func(c, r int) bool { return c >= 0 && r >= 0 } // a quarter of the plane
	if m := Mask(in, 0, 0); m != East|South|SouthEast {
		t.Errorf("the corner cell's mask is %08b", m)
	}
	if m := Mask(in, 5, 5); m != AllSides {
		t.Errorf("an inner cell's mask is %08b", m)
	}
	if Canonical(North|NorthEast) != North {
		t.Error("a corner beside land counted")
	}
}

func TestTheLastAreaOfALayerDecides(t *testing.T) {
	g := meadow(t)
	sawWater, sawTree := false, false
	for r := -60; r < 60; r++ {
		for c := -60; c < 60; c++ {
			if g.Tile(0, c, r) != grass {
				t.Fatal("the ground is not grass everywhere")
			}
			wet := g.Tile(1, c, r) == water
			sawWater = sawWater || wet
			if wet != g.In("water", c, r) {
				t.Fatal("water where the area is not, or none where it is")
			}
			if g.Tile(2, c, r) == tree {
				sawTree = true
				if g.Near("water", c, r, 1) {
					t.Fatalf("a tree at %d,%d beside the water", c, r)
				}
			}
		}
	}
	if !sawWater || !sawTree {
		t.Errorf("water %v, trees %v in 14400 cells", sawWater, sawTree)
	}
}

// Wherever the camera goes and whatever size the window, every pixel is
// drawn: chunks are made before they are seen, and the background never
// shows.
func TestTheWorldNeverLetsTheBackgroundShow(t *testing.T) {
	g := game.New()
	g.Background = canvas.RGB(255, 0, 255)
	w := New("world", meadow(t), 8, 0, 1, 2)
	g.Add(w.Object)
	sizes := [][2]int{{40, 30}, {120, 70}, {13, 90}}
	for i := range 120 {
		sz := sizes[i/40]
		cv, _ := canvas.NewCanvas(sz[0], sz[1])
		g.Draw(cv) // the view's size, as Run knows it from the window
		g.Camera = vector.Vector3{X: float64(i)*7.3 - 200, Y: math.Sin(float64(i)/9) * 150}
		g.Step(1.0 / 60)
		g.Draw(cv)
		for y := range cv.Height {
			for x := range cv.Width {
				if cv.At(x, y) == g.Background {
					t.Fatalf("frame %d, camera %v: the background shows at %d,%d", i, g.Camera, x, y)
				}
			}
		}
	}
	// Chunks left behind are let go.
	if n := len(w.chunks); n > 4*4*2 {
		t.Errorf("%d chunks kept", n)
	}
}

// A chunk shows the tiles the generator gives its cells, wherever its
// edges fall: trees and shores across a chunk's edge agree.
func TestChunksShowTheGeneratorsTiles(t *testing.T) {
	gen := meadow(t)
	g := game.New()
	w := New("world", gen, 8, 0, 1, 2)
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(64, 64)
	g.Draw(cv)
	g.Step(1.0 / 60)
	for _, ch := range w.chunks {
		for l, m := range ch.layers {
			for r := range 8 {
				for c := range 8 {
					if got, want := m.At(c, r), gen.Tile(l, ch.cx*8+c, ch.cy*8+r); got != want {
						t.Fatalf("chunk %d,%d layer %d cell %d,%d: %d, not %d", ch.cx, ch.cy, l, c, r, got, want)
					}
				}
			}
		}
	}
}

// The world's water is a wall to walk into.
func TestTheWorldsWaterIsAWall(t *testing.T) {
	gen := meadow(t)
	g := game.New()
	w := New("world", gen, 8, 0, 1, 2)
	g.Add(w.Object)
	cv, _ := canvas.NewCanvas(200, 200)
	// Find water, and stand left of it.
	col, row := 0, 0
	for !gen.In("water", col, row) || gen.In("water", col-1, row) {
		if col++; col > 200 {
			col, row = 0, row+1
		}
	}
	g.Camera = vector.Vector3{X: float64(col*4) - 100, Y: float64(row*4) - 100}
	g.Draw(cv)
	g.Step(1.0 / 60)
	o := game.NewObject("walker")
	o.Collider = collider.Box{X: -1, Y: -1, Width: 2, Height: 2}
	o.Transform.Position = vector.Vector3{X: float64(col*4) - 2, Y: float64(row*4) + 2}
	if hit, _ := tilemap.Move(o, 20, 0, w); !hit || o.Transform.Position.X > float64(col*4)-0.98 {
		t.Errorf("walked into the water at x %d to %v", col*4, o.Transform.Position.X)
	}
}

// A point of the world is in an area as the cell under it is, the
// negative side of the origin and all.
func TestInAtIsTheAreaOfTheCellUnderAPoint(t *testing.T) {
	gen := meadow(t)
	w := New("world", gen, 8, 0, 1, 2)
	for r := -20; r < 20; r++ {
		for c := -20; c < 20; c++ {
			want := gen.In("water", c, r)
			for _, d := range [][2]float64{{0.1, 0.1}, {3.9, 3.9}} {
				if got := w.InAt("water", float64(c*4)+d[0], float64(r*4)+d[1]); got != want {
					t.Fatalf("InAt in cell %d,%d is %v, not %v", c, r, got, want)
				}
			}
		}
	}
}
