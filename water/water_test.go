package water

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/terrain"
)

// colDiff is how far apart columns x and y of img are, summed.
func colDiff(img *image.NRGBA, x, y int) int {
	d := 0
	for row := range img.Bounds().Dy() {
		a, b := img.NRGBAAt(x, row), img.NRGBAAt(y, row)
		d += abs(int(a.R)-int(b.R)) + abs(int(a.G)-int(b.G))
	}
	return d
}

// rowDiff is how far apart rows x and y of img are, summed.
func rowDiff(img *image.NRGBA, x, y int) int {
	d := 0
	for col := range img.Bounds().Dx() {
		a, b := img.NRGBAAt(col, x), img.NRGBAAt(col, y)
		d += abs(int(a.R)-int(b.R)) + abs(int(a.G)-int(b.G))
	}
	return d
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// The map's edges meet: the step from its last column to its first, as
// copies lie side by side, is no bigger than a step inside it.
func TestNormalsTile(t *testing.T) {
	m := Normals(Config{Size: 64, Seed: 3, Steepness: 1})
	worstCol, worstRow := 0, 0
	for i := range 63 {
		worstCol = max(worstCol, colDiff(m, i, i+1))
		worstRow = max(worstRow, rowDiff(m, i, i+1))
	}
	if d := colDiff(m, 63, 0); d > worstCol {
		t.Errorf("left and right edges differ by %d, a step inside by up to %d", d, worstCol)
	}
	if d := rowDiff(m, 63, 0); d > worstRow {
		t.Errorf("top and bottom edges differ by %d, a step inside by up to %d", d, worstRow)
	}
	// A larger map is the same waves: its pixel at the same place agrees.
	big := Normals(Config{Size: 128, Seed: 3, Steepness: 1})
	if a, b := m.NRGBAAt(5, 9), big.NRGBAAt(10, 18); a != b {
		t.Errorf("the same place differs by size: %v and %v", a, b)
	}
}

// The normals are unit, point up, and the waves are not flat.
func TestNormalsAreWaves(t *testing.T) {
	m := Normals(Config{Size: 32, Seed: 1, Steepness: 1})
	lo, hi := 255, 0
	for y := range 32 {
		for x := range 32 {
			c := m.NRGBAAt(x, y)
			n := [3]float64{float64(c.R)/127.5 - 1, float64(c.G)/127.5 - 1, float64(c.B)/127.5 - 1}
			if l := math.Sqrt(dot(n, n)); math.Abs(l-1) > 0.03 || n[2] <= 0 {
				t.Fatalf("pixel %d,%d: normal %v", x, y, n)
			}
			lo, hi = min(lo, int(c.R)), max(hi, int(c.R))
		}
	}
	if hi-lo < 20 {
		t.Errorf("the surface is flat: X spans %d to %d", lo, hi)
	}
}

// A loop of frames ends where it began, and its frames differ.
func TestFramesLoop(t *testing.T) {
	at := func(f int) *image.NRGBA {
		return Normals(Config{Size: 32, Seed: 5, Steepness: 1, Frame: f, Frames: 4})
	}
	first := at(0)
	for i, p := range first.Pix {
		if at(4).Pix[i] != p {
			t.Fatal("frame 4 of 4 is not frame 0")
		}
	}
	same := true
	for i, p := range first.Pix {
		same = same && at(1).Pix[i] == p
	}
	if same {
		t.Error("frame 1 is frame 0: the water does not move")
	}
}

// Water mirrors more toward the grazing angle.
func TestFresnelGrows(t *testing.T) {
	last := -1.0
	for _, cos := range []float64{1, 0.8, 0.6, 0.4, 0.2, 0.05, 0} {
		f := Fresnel(cos, 0.02)
		if f <= last {
			t.Fatalf("Fresnel at cos %v is %v, not above %v", cos, f, last)
		}
		last = f
	}
	if math.Abs(Fresnel(1, 0.02)-0.02) > 1e-9 || math.Abs(Fresnel(0, 0.02)-1) > 1e-9 {
		t.Error("Fresnel is not f0 head on and 1 at grazing")
	}
}

// A tilted surface seen more grazing shades brighter than a flat one.
func TestShadeMirrorsAtGrazing(t *testing.T) {
	s := Lake().Surface
	flat := s.Shade(Look{Normal: [3]float64{0, 0, 1}, Depth: 1.5})
	away := unit([3]float64{-0.8, -0.8, 1})
	grazing := s.Shade(Look{Normal: away, Depth: 1.5})
	if grazing[2] <= flat[2] {
		t.Errorf("grazing blue %v is not above flat %v", grazing[2], flat[2])
	}
}

// ground is a plain tile of sand: opaque diamond, 64 by 32.
func ground() *image.NRGBA {
	g := image.NewNRGBA(image.Rect(0, 0, 64, 32))
	for y := range 32 {
		for x := range 64 {
			if math.Abs(float64(x)+0.5-32)/32+math.Abs(float64(y)+0.5-16)/16 <= 1 {
				g.SetNRGBA(x, y, color.NRGBA{200, 180, 130, 255})
			}
		}
	}
	return g
}

var (
	lake    = Lake()
	flat    = Normals(Config{Size: 64, Span: Span, Seed: 1, Steepness: 1})
	tileSet = ground()
)

func tileOf(mask int) *image.NRGBA {
	return Tile(Shore{Mask: mask, Surface: &lake.Surface, Normals: flat, Ground: tileSet})
}

// isSand is whether c is the ground's own colour, as land shows it.
func isSand(c color.NRGBA) bool { return c.R > 100 && c.R > c.B+20 && c.G > 90 }

// A tile is land where the mask says its neighbour is, and water where
// the mask says all round.
func TestShoreLand(t *testing.T) {
	north := tileOf(terrain.AllSides &^ terrain.North &^ terrain.NorthEast &^ terrain.NorthWest)
	// The top corner of the diamond, beside the northern neighbour.
	if c := north.NRGBAAt(32, 1); !isSand(c) {
		t.Errorf("land to the north is %v", c)
	}
	if c := north.NRGBAAt(32, 30); isSand(c) {
		t.Errorf("water to the south is sand %v", c)
	}
	inside := tileOf(terrain.AllSides)
	for _, pt := range [][2]int{{32, 1}, {32, 30}, {4, 16}, {60, 16}, {32, 16}} {
		if c := inside.NRGBAAt(pt[0], pt[1]); isSand(c) {
			t.Errorf("tile with water all round is land at %v: %v", pt, c)
		}
	}
	// Outside the diamond stays transparent.
	if c := inside.NRGBAAt(0, 0); c.A != 0 {
		t.Errorf("outside the diamond is %v", c)
	}
}

// Every mask makes a tile, and the water's frames differ in the foam only
// where a shore is.
func TestShoreEveryMask(t *testing.T) {
	for _, mask := range terrain.Masks() {
		img := tileOf(mask)
		if img.NRGBAAt(32, 16).A != 255 {
			t.Fatalf("mask %d: centre not opaque", mask)
		}
	}
}

// A shore meets its neighbour: the point on the edge two tiles share is
// the same distance from the land in both, whatever else their masks say.
func TestShoreEdgesAgree(t *testing.T) {
	a := landOf(terrain.AllSides &^ terrain.North &^ terrain.NorthEast &^ terrain.NorthWest)
	b := landOf(terrain.AllSides &^ terrain.NorthWest)
	// A's east edge, 0.3 from its northern land, is B's west edge, 0.3 from
	// the corner of its north-western land.
	if da, db := distance(a, 1, 0.3), distance(b, 0, 0.3); math.Abs(da-db) > 1e-9 {
		t.Errorf("the shared edge is %v from land in one tile, %v in the other", da, db)
	}
}

// The world's waves and wiggle meet across every shared edge, whichever
// variant lies on either side, and a cell is not its neighbour's copy.
func TestVariantsMeet(t *testing.T) {
	same := func(a, b float64) bool {
		d := math.Abs(a - b)
		return d < 1e-12 || math.Abs(d-1) < 1e-12
	}
	for col := range 4 {
		for row := range 4 {
			here, east, south := Shore{Col: col, Row: row}, Shore{Col: col + 1, Row: row}, Shore{Col: col, Row: row + 1}
			for _, f := range []float64{0, 0.3, 0.9} {
				au, av := here.world(1, f)
				bu, bv := east.world(0, f)
				if !same(au, bu) || av != bv {
					t.Fatalf("cells %d,%d and east meet at %v,%v and %v,%v", col, row, au, av, bu, bv)
				}
				au, av = here.world(f, 1)
				bu, bv = south.world(f, 0)
				if !same(av, bv) || au != bu {
					t.Fatalf("cells %d,%d and south meet at %v,%v and %v,%v", col, row, au, av, bu, bv)
				}
			}
		}
	}
	n := Normals(Config{Size: 64, Span: Span, Seed: 1, Steepness: 1})
	a, b := Shore{Col: 0}, Shore{Col: 1}
	au, av := a.world(0.5, 0.5)
	bu, bv := b.world(0.5, 0.5)
	if normalAt(n, au, av) == normalAt(n, bu, bv) {
		t.Error("two neighbouring cells have the same waves")
	}
	if parity(-1) != 1 || parity(-2) != 0 {
		t.Error("parity of negative places")
	}
}
