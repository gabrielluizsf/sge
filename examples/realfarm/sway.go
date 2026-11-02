package main

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// swayHeight is how far above its foot a plant leans the whole of its
// sway, in pixels.
const swayHeight = 40.0

// sway is the frame cv with whatever rises above the cell's middle leaned
// px pixels to the right at the height of swayHeight, more the higher it
// is and none at the foot: a plant bending in the wind.
func sway(cv *canvas.Canvas, px float64) *canvas.Canvas {
	return swayAt(cv, swing{px: px, height: swayHeight, foot: centreAt[1]})
}

// swing is a lean: px pixels at height pixels above the foot, which is
// foot rows down the picture.
type swing struct {
	px, height float64
	foot       int
}

// swayAt is cv with what rises above its foot leaned as s says; cv's size
// is the result's.
func swayAt(cv *canvas.Canvas, s swing) *canvas.Canvas {
	out, err := canvas.NewCanvas(cv.Width, cv.Height)
	if err != nil {
		panic(err)
	}
	for y := 0; y < cv.Height; y++ {
		rise := math.Max(0, float64(s.foot-y))
		dx := int(math.Round(s.px * rise / s.height))
		for x := 0; x < cv.Width; x++ {
			if nx := x + dx; nx >= 0 && nx < out.Width {
				out.Pixels[y*out.Stride+nx] = cv.Pixels[y*cv.Stride+x]
			}
		}
	}
	return out
}

// swayOf is the lean of frame k of the four a plant sways through.
func swayOf(k int) float64 { return [4]float64{0, 1.5, 0, -1.5}[k%4] }
