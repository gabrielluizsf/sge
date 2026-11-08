// Package vignette tints the edges of a screen with a colour that grows
// smoothly from nothing in the middle: a wound's red, a cold's blue, a
// dream's haze. It is made once for a size and drawn fast, for it covers
// much of the screen in a frame.
package vignette

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
)

// Config is what a vignette looks like: its Color at the very edge, how
// opaque there, from 0 to 255, and where it begins, From, from 0, the
// middle of the screen, to 1, its edge.
type Config struct {
	Color canvas.Color
	Alpha int
	From  float64
}

// Vignette is a Config worked out for a screen's size.
type Vignette struct {
	Width, Height int
	cfg           Config
	alpha         []uint8 // row by row
	reach         []int   // by row: how far in from each side it has any alpha
}

// New is the vignette cfg for a screen w by h: a rounded box, following
// the screen, more of it at the corners.
func New(w, h int, cfg Config) *Vignette {
	v := &Vignette{Width: w, Height: h, cfg: cfg, alpha: make([]uint8, w*h), reach: make([]int, h)}
	for y := range h {
		ny := (float64(y)+0.5)/float64(h)*2 - 1
		for x := range w {
			nx := (float64(x)+0.5)/float64(w)*2 - 1
			r := math.Pow(nx*nx*nx*nx+ny*ny*ny*ny, 0.25)
			t := min(1, max(0, (r-cfg.From)/(1-cfg.From)))
			if a := int(float64(cfg.Alpha)*t*t*(3-2*t) + 0.5); a > 0 {
				v.alpha[y*w+x] = uint8(a)
				if x < w/2 {
					v.reach[y] = max(v.reach[y], x+1)
				} else {
					v.reach[y] = max(v.reach[y], w-x)
				}
			}
		}
	}
	return v
}

// AlphaAt is how opaque the vignette is at x, y at full strength.
func (v *Vignette) AlphaAt(x, y int) uint8 { return v.alpha[y*v.Width+x] }
