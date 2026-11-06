package daytime

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// Shade is the script of an object fixed to the screen, over the world,
// that lights it as the clock's hour does: every pixel is multiplied by
// the hour's Ambient, so the picture keeps its contrast and its blacks,
// where a colour laid over it would wash it out. In full day it does
// nothing.
type Shade struct {
	Clock *Clock
}

func (s Shade) Update(*game.Object, float64) {}

func (s Shade) Paint(_ *game.Object, cv *canvas.Canvas) {
	a := s.Clock.Ambient()
	if a == (Ambient{1, 1, 1}) {
		return
	}
	if cv.Pixels == nil {
		multiplySlow(cv, a)
		return
	}
	var r, g, b [256]uint32
	for i := range r {
		v := float64(i)
		r[i] = uint32(v*a[0]+0.5) << 16
		g[i] = uint32(v*a[1]+0.5) << 8
		b[i] = uint32(v*a[2] + 0.5)
	}
	for y := range cv.Height {
		row := cv.Pixels[y*cv.Stride : y*cv.Stride+cv.Width]
		for x, p := range row {
			row[x] = canvas.Color(uint32(p)&0xFF000000 | r[uint8(p>>16)] | g[uint8(p>>8)] | b[uint8(p)])
		}
	}
}

// multiplySlow lights a canvas whose pixels cannot be reached directly.
func multiplySlow(cv *canvas.Canvas, a Ambient) {
	for y := range cv.Height {
		for x := range cv.Width {
			p := cv.At(x, y)
			cv.Put(x, y, canvas.RGBA(canvas.RGBAParams{
				R: uint8(float64(p.R())*a[0] + 0.5), G: uint8(float64(p.G())*a[1] + 0.5), B: uint8(float64(p.B())*a[2] + 0.5), A: p.A(),
			}))
		}
	}
}
