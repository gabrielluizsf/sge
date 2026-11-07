package main

import (
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/face"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/sprite"
)

// faceName is the baked face of how a face looks: a blink before the
// mouth, which a blink's moment hides; none for the still face.
func faceName(look face.Look) string {
	switch look.Lids {
	case face.LidsHalf:
		return "lids1"
	case face.LidsShut:
		return "lids2"
	}
	switch look.Mouth {
	case face.Ajar:
		return "say1"
	case face.Open:
		return "say2"
	case face.Round:
		return "say3"
	}
	return ""
}

// Face is who's portrait with their face as look has it, baked at their
// time of life or the nearest; their still portrait for the still face or
// while the face is not baked.
func (a *realArt) Face(who farm.Who, look face.Look) *canvas.Canvas {
	f := faceName(look)
	if f == "" {
		return a.Portrait(who)
	}
	id := villagerOf(who)
	for _, stage := range stageOrder[who.Stage] {
		if !a.shelf.Has(bakes.PortraitName(id, stage)) {
			continue
		}
		name := bakes.PortraitFaceName(id, stage, f)
		if cv, ok := a.hud["p"+name]; ok {
			return cv
		}
		if !a.shelf.Has(name) {
			return a.Portrait(who)
		}
		s, err := sprite.Load(filepath.Join(a.shelf.Dir, name+".png"))
		if err != nil {
			return a.Portrait(who)
		}
		cv := shrink(s.Image, portraitW, portraitH)
		a.hud["p"+name] = cv
		return cv
	}
	return a.Portrait(who)
}
