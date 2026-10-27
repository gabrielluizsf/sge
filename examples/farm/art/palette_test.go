package art

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/terrain"
)

// size is the drawing's width and height, as its sprite's kind says.
func size(name string) (w, h int) {
	switch {
	case strings.HasPrefix(name, "farmer_"):
		return FarmerW, FarmerH
	case strings.HasPrefix(name, "icon_"):
		return Size, Size
	case name == "cursor":
		return IsoW, Size
	case name == "stall":
		return StallW, StallH
	}
	return IsoW, IsoH
}

// render is the sprite name drawn as the game draws it, twice the size.
func render(t *testing.T, name, src string) *canvas.Canvas {
	w, h := size(name)
	s, err := sprite.FromSVG(src, 2*w, 2*h)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return s.Image
}

// Every sprite is SVG the game can draw, of its kind's size.
func TestSpritesAreDrawingsOfTheirSize(t *testing.T) {
	for name, src := range Sprites() {
		w, h := size(name)
		if !strings.Contains(src, fmt.Sprintf(`viewBox="0 0 %d %d"`, w, h)) {
			t.Errorf("%s is not %d by %d", name, w, h)
		}
		render(t, name, src)
	}
}

// inDiamond reports whether pixel x, y of a tile drawn twice the size has
// its middle on its cell's diamond.
func inDiamond(x, y int) bool {
	sx, sy := (float64(x)+0.5)/2-Size, (float64(y)+0.5)/2-2*Size
	u, v := sy+sx/2, sy-sx/2
	return u > 0.05 && v > 0.05 && u < Size-0.05 && v < Size-0.05
}

// outDiamond reports whether pixel x, y lies clear of the diamond.
func outDiamond(x, y int) bool {
	sx, sy := (float64(x)+0.5)/2-Size, (float64(y)+0.5)/2-2*Size
	u, v := sy+sx/2, sy-sx/2
	return u < -0.3 || v < -0.3 || u > Size+0.3 || v > Size+0.3
}

// Ground tiles, and open water, cover their cell's diamond whole, so tiles
// side by side meet without a seam, and nothing beyond it.
func TestGroundCoversItsDiamond(t *testing.T) {
	ground := []string{"grass0", "grass1", "flowers2", "path", "soil_dry", "soil_wet", fmt.Sprint("shore", terrain.AllSides, "_1")}
	all := Sprites()
	for _, name := range ground {
		img := render(t, name, all[name])
		holes, spills := 0, 0
		for y := range img.Height {
			for x := range img.Width {
				a := img.At(x, y).A()
				if inDiamond(x, y) && a != 255 {
					holes++
				}
				if outDiamond(x, y) && a != 0 {
					spills++
				}
			}
		}
		if holes > 0 || spills > 0 {
			t.Errorf("%s: %d pixels of its diamond not opaque, %d drawn beyond it", name, holes, spills)
		}
	}
}

// A pond's tile ends in a bank on its sides facing land: transparent at
// those edges' middles, so the grass under the water shows.
func TestShoresOpenOntoTheLand(t *testing.T) {
	img := render(t, "shore0_0", shoreTile(0, 0))
	// Near the middles of the diamond's four edges, just inside.
	for _, p := range [][2]int{{48, 72}, {48, 88}, {16, 88}, {16, 72}} {
		if a := img.At(p[0], p[1]).A(); a != 0 {
			t.Errorf("an island's tile is not land at %v: alpha %d", p, a)
		}
	}
	if a := img.At(32, 80).A(); a != 255 {
		t.Errorf("an island's tile has no water in its middle")
	}
}

// Objects stand on the ground rather than carry a patch of it.
func TestObjectsAreTransparentInTheirCorners(t *testing.T) {
	for name, src := range Sprites() {
		if name != "stall" && name != "fence" && !strings.HasPrefix(name, "tree") && !strings.HasPrefix(name, "farmer_") {
			continue
		}
		img := render(t, name, src)
		for _, p := range [][2]int{{0, 0}, {img.Width - 1, 0}, {0, img.Height - 1}, {img.Width - 1, img.Height - 1}} {
			if img.At(p[0], p[1]).A() != 0 {
				t.Errorf("%s paints its corner %v", name, p)
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
