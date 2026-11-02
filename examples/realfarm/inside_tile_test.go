package main

import (
	"encoding/json"
	"image"
	"slices"
	"testing"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// No tile of the realistic farm is the stylized one's, the houses'
// insides, the forge and the stumps among them: before the baker has made
// a tile it is blank.
func TestNoTileIsTheStylizedArts(t *testing.T) {
	real, style := testArt(t), farm.SVGArt()
	for _, name := range farm.TileNames() {
		got := real.Tile(name)
		for i, p := range got.Pixels {
			if p.A() != 0 {
				t.Errorf("%s draws pixel %d with nothing baked", name, i)
				break
			}
		}
	}
	if b := real.Building("house0"); b == nil || same(b.Image, style.Building("house0").Image) {
		t.Error("an unbaked house is the stylized one")
	}
}

// Every piece of a house's inside the farm lays is baked: each floor,
// rug and back wall of every kind of house by package building, each
// piece of furniture from a model of the manifest.
func TestEveryInsidePieceIsBaked(t *testing.T) {
	names := farm.TileNames()
	pieces := map[string]bool{}
	for _, k := range building.Insides() {
		pieces[k.Name] = true
		if !isInside(k.Name) || !slices.Contains(names, k.Name) {
			t.Errorf("piece %s is no tile the farm lays", k.Name)
		}
	}
	if art.WallKinds < len(building.Houses) {
		t.Errorf("walls for %d kinds of house, of %d", art.WallKinds, len(building.Houses))
	}
	for k := range len(building.Houses) {
		for _, w := range []string{"wall_n", "wall_nwin", "wall_w", "wall_nw"} {
			if name := w + string(rune('0'+k)); !pieces[name] {
				t.Errorf("no piece %s", name)
			}
		}
	}
	var m kit.Manifest
	if err := json.Unmarshal(assetsJSON, &m); err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, a := range m.Assets {
		ids[a.ID] = a.Kind == "model"
	}
	for name, f := range furnishings {
		if !slices.Contains(names, name) || !ids[f.id] {
			t.Errorf("%s: %s is not a model of the manifest, or no tile", name, f.id)
		}
	}
	if _, ok := testArt(t).(farm.FurnitureArt); !ok {
		t.Fatal("the realistic art does not say how large its furniture is")
	}
	if c, r := testArt(t).(farm.FurnitureArt).Furniture("bed"); c*r != 2 {
		t.Errorf("a bed covers %d by %d cells, not two", c, r)
	}
	if bakes.InsideName("floor") == "floor" {
		t.Error("an inside's sprite shares the tile's name")
	}
}

// A bed, drawn from its back cell, has its picture's middle half a cell
// forward, down and right on the screen, between its two cells.
func TestABedReachesForwardFromItsBackCell(t *testing.T) {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	img := image.NewNRGBA(image.Rect(0, 0, 40, 20))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	if err := shelf.Put(furnishings["bed"].id+"_0", bake.Sprite{Image: img, Pivot: image.Pt(20, 10)}); err != nil {
		t.Fatal(err)
	}
	a := newArt(shelf).(*realArt)
	cv := a.Tile("bed")
	// The pivot stands half a cell, TileW/4 across and TileW/8 down, from
	// the diamond's middle at the bottom of the frame.
	left, top := cv.Width, cv.Height
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y).A() > 0 {
				left, top = min(left, x), min(top, y)
			}
		}
	}
	wantX := cv.Width/2 + farm.TileW/4 - 20
	wantY := cv.Height - farm.TileW/4 + farm.TileW/8 - 10
	if left != wantX || top != wantY {
		t.Errorf("the bed's picture starts at %d,%d, not %d,%d", left, top, wantX, wantY)
	}
}
