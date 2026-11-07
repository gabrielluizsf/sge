package main

import (
	"image/color"
	"testing"
)

// centre is the stage picture's middle pixel.
func centre(t *testing.T, s *liveStage, v *liveView) color.NRGBA {
	t.Helper()
	img, err := s.Render(v, 64, 80)
	if err != nil || img == nil {
		t.Fatalf("render: %v, %v", img, err)
	}
	c := img.At(32, 40)
	return color.NRGBA{c.R(), c.G(), c.B(), 255}
}

// The person is drawn live, and a recolour shows in the next picture
// without the model being uploaded again.
func TestTheStageDrawsThePersonAndRecolours(t *testing.T) {
	s := newLiveStage(softDevice)
	defer s.Close()
	v := &liveView{}
	if img, _ := s.Render(v, 64, 80); img != nil {
		t.Error("a picture with nobody on the stage")
	}
	f := &fakeMaker{}
	p, _ := f.Make(someone(1))
	if err := s.Show(&liveUpdate{person: p, fresh: true}, 64, 80); err != nil {
		t.Fatal(err)
	}
	dark := centre(t, s, v)
	model := s.model
	apply := p.(*fakePerson).Plan(withTone(someone(1), 1))
	if err := s.Show(&liveUpdate{person: p, applies: []func() []int{apply}}, 64, 80); err != nil {
		t.Fatal(err)
	}
	light := centre(t, s, v)
	if s.model != model {
		t.Error("the model was uploaded again for a colour")
	}
	if light.R <= dark.R+60 {
		t.Errorf("centre %v, then %v: the colour did not change", dark, light)
	}
	// Nothing changed: the same picture is given without drawing it again.
	img1, _ := s.Render(v, 64, 80)
	img2, _ := s.Render(v, 64, 80)
	if img1 != img2 {
		t.Error("drew again with nothing changed")
	}
	v.Turn(30)
	s.Dirty()
	if _, err := s.Render(v, 64, 80); err != nil || s.dirty {
		t.Errorf("render after a turn: %v", err)
	}
}
