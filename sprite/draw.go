package sprite

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
)

// Draw paints the sprite onto dst where t places it.
func (s *Sprite) Draw(dst *canvas.Canvas, t transform.Transform) {
	if s == nil || s.Image == nil {
		return
	}
	area := canvas.Area{Width: s.Image.Width, Height: s.Image.Height}
	dst.BlitMatrix(s.Image, s.Matrix(t), area)
}
