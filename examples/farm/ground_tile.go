package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// The grounds the tiles lie on.
const (
	grassGround = "sparse_grass"
	soilGround  = "farm_soil"
	pathGround  = "grass_path_2"
)

// groundTile is the frame of patch k of ground id: the cell's diamond at
// the bottom, nothing above it. Nil while the ground is not baked.
func (a *realArt) groundTile(id string, k int) *canvas.Canvas {
	g, ok := a.bake(fmt.Sprintf("%s_g%d", id, k))
	if !ok {
		return nil
	}
	cv := newFrame()
	g.put(cv, cornerAt)
	return cv
}

// flowersTile is patch k of grass with dandelions on it, swaying with k.
func (a *realArt) flowersTile(k int) *canvas.Canvas {
	cv := a.groundTile(grassGround, k)
	if cv == nil {
		return nil
	}
	if f, ok := a.first("dandelion_01_0", "dandelion_01_v0_0"); ok {
		top := newFrame()
		f.put(top, centreAt)
		cv.BlitOver(0, 0, sway(top, swayOf(k)))
	}
	return cv
}

// wetSoilTile is the tilled soil darkened by water.
func (a *realArt) wetSoilTile() *canvas.Canvas {
	cv := a.groundTile(soilGround, 0)
	if cv == nil {
		return nil
	}
	for i, c := range cv.Pixels {
		if c>>24 != 0 {
			cv.Pixels[i] = c&0xFF000000 | (c>>16&0xFF)*55/100<<16 | (c>>8&0xFF)*55/100<<8 | (c&0xFF)*62/100
		}
	}
	return cv
}
