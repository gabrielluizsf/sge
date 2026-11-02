package main

import (
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/sprite"
)

// The room a portrait has in the conversation box's panel, in pixels.
const (
	portraitW = 52 * farm.Scale
	portraitH = 72 * farm.Scale
)

// Portrait is who's head and shoulders, baked from the same model as the
// person who walks the farm, at their time of life or the nearest baked;
// nil while none is.
func (a *realArt) Portrait(who farm.Who) *canvas.Canvas {
	id := villagerOf(who)
	for _, stage := range stageOrder[who.Stage] {
		name := bakes.PortraitName(id, stage)
		if cv, ok := a.hud["p"+name]; ok {
			return cv
		}
		if !a.shelf.Has(name) {
			continue
		}
		s, err := sprite.Load(filepath.Join(a.shelf.Dir, name+".png"))
		if err != nil {
			continue
		}
		cv := shrink(s.Image, portraitW, portraitH)
		a.hud["p"+name] = cv
		return cv
	}
	return nil
}
