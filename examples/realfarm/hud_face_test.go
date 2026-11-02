package main

import (
	"image"
	"image/color"
	"testing"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/face"
	"github.com/gabrielluizsf/sge/life"
)

// A talking villager's portrait is their baked face as it looks: a blink
// before the mouth, the still portrait for the still face, and the still
// one too while the faces are not baked.
func TestTheTalkingPortraitIsTheFaceAsItLooks(t *testing.T) {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	a := newArt(shelf).(*realArt)
	who := farm.Who{Stage: life.Adult}
	if a.Face(who, face.Look{Lids: face.LidsShut}) != nil {
		t.Fatal("a face with nothing baked")
	}
	id := villagerOf(who)
	paint := func(name string, shade uint8) {
		img := image.NewNRGBA(image.Rect(0, 0, 40, 50))
		for i := range img.Pix {
			img.Pix[i] = shade
		}
		img.SetNRGBA(0, 0, color.NRGBA{shade, shade, shade, 255})
		if err := shelf.Put(name, bake.Sprite{Image: img, Pivot: image.Pt(20, 50)}); err != nil {
			t.Fatal(err)
		}
	}
	paint(bakes.PortraitName(id, life.Adult), 10)
	if same(a.Face(who, face.Look{Mouth: face.Open}), a.Portrait(who)) == false {
		t.Error("an unbaked face is not the still portrait")
	}
	for i, f := range bakes.PortraitFaces {
		paint(bakes.PortraitFaceName(id, life.Adult, f.Name), uint8(20+i*20))
	}
	a = newArt(shelf).(*realArt)
	still := a.Portrait(who)
	looks := map[string]face.Look{
		"lids1": {Lids: face.LidsHalf, Mouth: face.Open}, "lids2": {Lids: face.LidsShut},
		"say1": {Mouth: face.Ajar}, "say2": {Mouth: face.Open}, "say3": {Mouth: face.Round},
	}
	for name, look := range looks {
		if faceName(look) != name {
			t.Errorf("%+v is %q, not %s", look, faceName(look), name)
		}
		if same(a.Face(who, look), still) {
			t.Errorf("%s is the still portrait", name)
		}
	}
	if !same(a.Face(who, face.Look{}), still) {
		t.Error("the still face is not the still portrait")
	}
	if _, ok := any(a).(farm.FaceArt); !ok {
		t.Error("the realistic art draws no faces")
	}
}
