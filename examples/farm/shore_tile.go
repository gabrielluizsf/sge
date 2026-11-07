package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// shoreTile is the water with water beside it as mask says, f steps into
// its ripple, the lake's, in the first of its variants.
func (a *realArt) shoreTile(mask, f int) *canvas.Canvas {
	b, ok := a.bake(fmt.Sprintf("water_lake_shore%d_00_%d", mask, f))
	if !ok {
		return nil
	}
	cv := newFrame()
	b.put(cv, cornerAt)
	return cv
}
