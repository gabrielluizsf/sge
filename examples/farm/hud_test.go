package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
)

// testArt is the realistic art over a shelf with nothing baked on it.
func testArt(t *testing.T) farm.Art {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return newArt(shelf)
}

// The realistic art draws a picture for every tile, crops at every stage
// and the withered plant too, and every icon of the crops, baked or not.
func TestTheRealArtDrawsEveryCrop(t *testing.T) {
	a := testArt(t)
	for _, name := range farm.TileNames() {
		if cv := a.Tile(name); cv == nil || cv.Width == 0 {
			t.Errorf("no picture for tile %s", name)
		}
	}
	for _, name := range []string{"icon_hoe", "icon_coffeeSeeds", "icon_sugarcane", "icon_cassava"} {
		if cv := a.Icon(name); cv == nil || cv.Width == 0 {
			t.Errorf("no picture for %s", name)
		}
	}
}

// Its text is Source Serif, accents and all, and its panels are the
// dark ones, translucent, the selected warmer than the rest.
func TestTheRealHUDIsSerifOnGlass(t *testing.T) {
	a := testArt(t)
	plain, accents := a.Text("acao", 'k'), a.Text("ação é só", 'k')
	if plain.Width == 0 || accents.Width <= plain.Width || a.Text("regador", 'w').Height == 0 {
		t.Fatalf("text %dx%d, %dx%d", plain.Width, plain.Height, accents.Width, accents.Height)
	}
	panel, chosen := a.Panel(40, 20, false), a.Panel(40, 20, true)
	p, s := panel.At(panel.Width/2, panel.Height/2), chosen.At(chosen.Width/2, chosen.Height/2)
	if p.A() == 0 || p.A() == 255 || p.R() > 80 || p.G() > 80 {
		t.Errorf("panel is not dark and translucent: %v", p)
	}
	if s.R() <= p.R() {
		t.Errorf("the chosen panel %v is not warmer than %v", s, p)
	}
	if a.Text("x", 'k').At(0, 0).A() != 0 {
		t.Error("text is not transparent round the glyph")
	}
}
