package main

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
)

// swayHeight is how far above its foot a plant leans the whole of its
// sway, in pixels.
const swayHeight = 40.0

// sway is the frame cv with whatever rises above the cell's middle leaned
// px pixels to the right at the height of swayHeight, more the higher it
// is and none at the foot: a plant bending in the wind.
func sway(cv *canvas.Canvas, px float64) *canvas.Canvas {
	out := newFrame()
	for y := 0; y < cv.Height; y++ {
		rise := math.Max(0, float64(centreAt[1]-y))
		dx := int(math.Round(px * rise / swayHeight))
		for x := 0; x < cv.Width; x++ {
			if nx := x + dx; nx >= 0 && nx < out.Width && y < farm.TileH {
				out.Pixels[y*out.Stride+nx] = cv.Pixels[y*cv.Stride+x]
			}
		}
	}
	return out
}

// swayOf is the lean of frame k of the four a plant sways through.
func swayOf(k int) float64 { return [4]float64{0, 1.5, 0, -1.5}[k%4] }
