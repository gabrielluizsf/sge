package sprite

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
)

// Draw paints the sprite onto dst where t places it. A sprite that is only
// moved, neither turned nor scaled, is copied whole onto the pixels the
// mapping would pick, rather than mapped back pixel by pixel.
func (s *Sprite) Draw(dst *canvas.Canvas, t transform.Transform) {
	if s == nil {
		return
	}
	s.ready()
	if s.Image == nil {
		return
	}
	m := s.Matrix(t)
	if m.A == 1 && m.D == 1 && m.B == 0 && m.C == 0 {
		// The mapping takes destination pixel x to source pixel
		// floor(x + 0.5 - E): the source sits at -floor(0.5 - E).
		dst.Blit(-int(math.Floor(0.5-m.E)), -int(math.Floor(0.5-m.F)), s.Image)
		return
	}
	area := canvas.Area{Width: s.Image.Width, Height: s.Image.Height}
	dst.BlitMatrix(s.Image, m, area)
}
