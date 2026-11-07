package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/farm"
)

// Where a baked sprite's pivot goes in a tile's frame: a ground's on the
// top corner of the cell's diamond, which sits at the bottom of the
// frame, and what stands on the cell on the diamond's middle.
var (
	cornerAt = [2]int{farm.TileW / 2, farm.TileH - farm.TileW/2}
	centreAt = [2]int{farm.TileW / 2, farm.TileH - farm.TileW/4}
)

// newFrame is a transparent picture of a tile's size.
func newFrame() *canvas.Canvas {
	cv, err := canvas.NewCanvas(farm.TileW, farm.TileH)
	if err != nil {
		panic(err)
	}
	return cv
}

// put draws b onto cv with its pivot at the point at.
func (b *baked) put(cv *canvas.Canvas, at [2]int) {
	cv.BlitOver(at[0]-b.at.PivotX, at[1]-b.at.PivotY, b.cv)
}
