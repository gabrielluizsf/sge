package sprite

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
)

// Matrix places the sprite's pixels on a canvas for transform t: around
// the pivot it is scaled, then turned, then moved to t's position.
func (s *Sprite) Matrix(t transform.Transform) canvas.Matrix {
	s.ready()
	px := s.PivotX * float64(s.Image.Width)
	py := s.PivotY * float64(s.Image.Height)
	return canvas.Translate(t.Position.X, t.Position.Y).
		Mul(canvas.Rotate(t.Angle())).
		Mul(canvas.Scale(t.Scale.X, t.Scale.Y)).
		Mul(canvas.Translate(-px, -py))
}
