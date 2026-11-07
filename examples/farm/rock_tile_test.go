package main

import (
	"image"
	"testing"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
)

// A stone in the bag is a baked rock when there is one, and nothing drawn
// otherwise: the realistic game draws nothing of the stylized look's.
func TestTheStoneIconIsBakedOrNothing(t *testing.T) {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	a := newArt(shelf).(*realArt)
	if a.iconOf("icon_stone") != nil {
		t.Error("an unbaked stone has a picture")
	}
	img := image.NewNRGBA(image.Rect(0, 0, bakes.IconSize, bakes.IconSize))
	if err := shelf.Put("icon_stone", bake.Sprite{Image: img}); err != nil {
		t.Fatal(err)
	}
	if a.iconOf("icon_stone") == nil {
		t.Error("a baked stone has no picture")
	}
}

// A rock wider than its cell keeps all of its picture: the part below the
// cell's bottom corner used to be cut off in a straight line.
func TestARockKeepsAllOfItsPicture(t *testing.T) {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	const w, h = 100, 70
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	if err := shelf.Put(rockModels[0]+"_0", bake.Sprite{Image: img, Pivot: image.Pt(40, 25)}); err != nil {
		t.Fatal(err)
	}
	a := newArt(shelf).(*realArt)
	for _, name := range farm.TileNames() {
		v, ok := farm.RockVariant(name)
		if !ok || v != 0 {
			continue
		}
		cv := a.rockTile(name)
		if cv == nil {
			continue
		}
		n := 0
		for y := range cv.Height {
			for x := range cv.Width {
				if cv.At(x, y).A() > 0 {
					n++
				}
			}
		}
		if n != w*h {
			t.Errorf("%s keeps %d of %d pixels", name, n, w*h)
		}
		return
	}
	t.Fatal("no rock was drawn")
}
