package art

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/terrain"
)

// Every sprite is the farm's: drawn in the palette alone, every row as
// wide as the first, the tiles IsoW by IsoH, the farmer FarmerW by
// FarmerH and the icons Size square.
func TestSpritesKeepToThePaletteAndSize(t *testing.T) {
	for name, g := range Sprites() {
		w, h := IsoW, IsoH
		switch {
		case strings.HasPrefix(name, "farmer_"):
			w, h = FarmerW, FarmerH
		case strings.HasPrefix(name, "icon_"):
			w, h = Size, Size
		case name == "cursor":
			w, h = IsoW, Size
		case name == "stall":
			w, h = len([]rune(g[0])), len(g)
		}
		if len(g) != h {
			t.Errorf("%s has %d rows, not %d", name, len(g), h)
		}
		for y, row := range g {
			if n := len([]rune(row)); n != w {
				t.Errorf("%s row %d is %d wide, not %d", name, y, n, w)
			}
			for x, c := range row {
				if _, ok := Palette[c]; !ok && c != '.' {
					t.Errorf("%s has %q at %d,%d, not in the palette", name, c, x, y)
				}
			}
		}
	}
}

// Ground tiles cover their cell's diamond whole, so tiles side by side
// meet without a seam, and nothing beyond it.
func TestGroundTilesCoverTheirDiamond(t *testing.T) {
	for name, g := range Sprites() {
		if !strings.HasPrefix(name, "grass") && !strings.HasPrefix(name, "flowers") && !strings.HasPrefix(name, "soil") && name != "path" {
			continue
		}
		for y, row := range g {
			for x, c := range []rune(row) {
				sx, sy := float64(x)+0.5-Size, float64(y-isoTop)+0.5
				u, v := sy+sx/2, sy-sx/2
				in := u >= 0 && v >= 0 && u < Size && v < Size
				if in && c == '.' {
					t.Errorf("%s has a hole at %d,%d", name, x, y)
				}
				if !in && c != '.' {
					t.Errorf("%s draws outside its diamond at %d,%d", name, x, y)
				}
			}
		}
	}
}

// Laid on its diamond, a square tile's far corner is the diamond's top
// and its near corner the bottom.
func TestIsoGroundLaysTheSquareOnTheDiamond(t *testing.T) {
	sq := blank(Size, Size).set(0, 0, 'r').set(Size-1, Size-1, 'y').set(Size-1, 0, 'w').set(0, Size-1, 'k')
	g := isoGround(sq)
	at := func(x, y int) rune { return []rune(g[y])[x] }
	if at(Size-1, isoTop) != 'r' || at(Size, isoTop) != 'r' {
		t.Errorf("the top of the diamond is not the square's far corner")
	}
	if at(Size-1, IsoH-1) != 'y' {
		t.Errorf("the bottom of the diamond is not the square's near corner")
	}
	// A square whose east half is white lays it on the diamond's right.
	half := blank(Size, Size)
	for y := range Size {
		half[y] = strings.Repeat("k", Size/2) + strings.Repeat("w", Size/2)
	}
	h := isoGround(half)
	if c := []rune(h[isoTop+Size/2])[IsoW-3]; c != 'w' {
		t.Errorf("the diamond's right is %q, not the square's east", c)
	}
	if c := []rune(h[isoTop+Size/2])[2]; c != 'k' {
		t.Errorf("the diamond's left is %q, not the square's west", c)
	}
}

func TestFontGlyphsAreFiveBySeven(t *testing.T) {
	for r, g := range font {
		for _, row := range g {
			if len(row) != 5 || strings.Trim(row, ".#") != "" {
				t.Errorf("glyph %q: %q", r, row)
			}
		}
	}
}

// SVG draws every pixel of the grid, as runs, in the palette's colours.
func TestSVGDrawsEveryInkedPixel(t *testing.T) {
	svg := SVG(grid(strings.Repeat("kk..yyyyyyyyyyyy\n", 16)))
	if strings.Count(svg, "<rect") != 32 || !strings.Contains(svg, `width="12"`) || !strings.Contains(svg, Palette['y']) {
		t.Errorf("SVG:\n%s", svg)
	}
}

// Objects stand on the ground rather than carry a patch of it: around
// them they are transparent, so the animated grass behind shows through.
func TestObjectsAreTransparentAroundThem(t *testing.T) {
	for name, g := range Sprites() {
		if name != "stall" && !strings.HasPrefix(name, "tree") && name != "fence" {
			continue
		}
		last := []rune(g[len(g)-1])
		if last[0] != '.' || last[len(last)-1] != '.' || []rune(g[0])[0] != '.' {
			t.Errorf("%s is not transparent in its corners: it paints the ground under it", name)
		}
	}
}

// A shore tile is open water on its watery sides, so it meets the water
// beside it, and ends in transparency on its sides facing land.
func TestShoresMeetTheWaterAndOpenOntoTheLand(t *testing.T) {
	for _, mask := range terrain.Masks() {
		g := shore(mask, 0)
		for _, s := range []struct {
			bit  int
			edge string
		}{{terrain.North, g[0]}, {terrain.South, g[Size-1]}, {terrain.West, column(g, 0)}, {terrain.East, column(g, Size-1)}} {
			mid := s.edge[5:11]
			if s.bit&mask != 0 && strings.ContainsAny(mid, ".k") {
				t.Errorf("shore %04b: side %d faces water but its edge is %q", mask, s.bit, s.edge)
			}
			if s.bit&mask == 0 && strings.Trim(s.edge, ".") != "" {
				t.Errorf("shore %04b: side %d faces land but its edge is %q", mask, s.bit, s.edge)
			}
		}
	}
	full := shore(terrain.AllSides, 2)
	for y, row := range water(2) {
		if full[y] != row {
			t.Fatal("open water is not the water tile")
		}
	}
}

func column(g Grid, x int) string {
	var b strings.Builder
	for _, row := range g {
		b.WriteByte(row[x])
	}
	return b.String()
}
