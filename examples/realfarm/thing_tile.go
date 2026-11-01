package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// treeTile is the tree, a sapling of the farm's species, leaning as frame
// k of its sway says.
func (a *realArt) treeTile(k int) *canvas.Canvas {
	b, ok := a.bake("tree_" + bakes.FarmTree + "_s0_0_0")
	if !ok {
		return nil
	}
	cv := newFrame()
	b.put(cv, centreAt)
	return sway(cv, swayOf(k))
}

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
	cv := newFrame()
	b.put(cv, centreAt)
	return cv
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
