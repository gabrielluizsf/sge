package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// thingTile is the model id standing on the middle of the cell, at its
// first heading; its first variant when the model is of several.
func (a *realArt) thingTile(id string) *canvas.Canvas {
	b, ok := a.first(id+"_0", id+"_v0_0")
	if !ok {
		return nil
	}
	cv := newFrame()
	b.put(cv, centreAt)
	return cv
}

// buildingTile is the building id baked at its first heading, standing on
// the middle of the cell.
func (a *realArt) buildingTile(id string) *canvas.Canvas {
	b, ok := a.bake(fmt.Sprintf("building_%s_0", id))
	if !ok {
		return nil
	}
	return b.standing() // a forge's chimney stands taller than a tile
}

// cropStages are the stage of the crop that grows each of the farm's
// five: the last two are both ripe, which sways.
var cropStages = [5]int{0, 1, 2, 4, 4}

// cropTile is the crop in stage k of the farm's five.
func (a *realArt) cropTile(crop string, k int) *canvas.Canvas {
	if k < 0 || k >= len(cropStages) {
		return nil
	}
	b, ok := a.bake(fmt.Sprintf("crop_%s_s%d_0", crop, cropStages[k]))
	if !ok {
		return nil
	}
	cv := newFrame()
	b.put(cv, centreAt)
	if k == len(cropStages)-1 {
		return sway(cv, 1.5)
	}
	return cv
}

// witheredTile is the plant that went a day dry, standing on the cell.
func (a *realArt) witheredTile() *canvas.Canvas {
	b, ok := a.bake("crop_withered_s0_0")
	if !ok {
		return nil
	}
	cv := newFrame()
	b.put(cv, centreAt)
	return cv
}
