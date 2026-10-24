package art

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/terrain"
)

// Every sprite is the farm's: Size by Size, drawn in the palette alone.
func TestSpritesKeepToThePaletteAndSize(t *testing.T) {
	for name, g := range Sprites() {
		if len(g) != Size {
			t.Errorf("%s has %d rows", name, len(g))
		}
		for y, row := range g {
			if n := len([]rune(row)); n != Size {
				t.Errorf("%s row %d is %d wide", name, y, n)
			}
			for x, c := range row {
				if _, ok := Palette[c]; !ok && c != '.' {
					t.Errorf("%s has %q at %d,%d, not in the palette", name, c, x, y)
				}
			}
		}
	}
}

// Tiles that repeat across the ground meet without a seam: their left and
// right edges, and top and bottom, could sit beside each other.
func TestGroundTilesHaveNoTransparentPixels(t *testing.T) {
	for name, g := range Sprites() {
		if !strings.HasPrefix(name, "grass") && !strings.HasPrefix(name, "water") && !strings.HasPrefix(name, "soil") && name != "path" {
			continue
		}
		for y, row := range g {
			if strings.ContainsRune(row, '.') {
				t.Errorf("%s has a hole in row %d", name, y)
			}
		}
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
		if !strings.HasPrefix(name, "stall") && !strings.HasPrefix(name, "tree") && name != "fence" {
			continue
		}
		holes := 0
		for _, row := range g {
			holes += strings.Count(row, ".")
		}
		if holes == 0 {
			t.Errorf("%s has no transparent pixels: it paints the ground under it", name)
		}
		if strings.ContainsAny(g[Size-1][:1]+g[Size-1][Size-1:], "123") {
			t.Errorf("%s has grass in its corners", name)
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
