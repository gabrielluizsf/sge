package game

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/transform"
)

// fadeObject draws o's look faded at t, and says so, when o covers a
// Reveal standing behind it; false when it is to be drawn as it is. Only
// a look that is moved, neither turned nor scaled, is faded.
func (g *Game) fadeObject(o *Object, t transform.Transform, cv *canvas.Canvas) bool {
	if len(g.covers) == 0 {
		return false
	}
	s := o.Look()
	if s == nil {
		return false
	}
	if s.Image == nil && s.Make == nil {
		return false
	}
	m := s.Matrix(t) // makes the sprite ready
	if s.Image == nil {
		return false
	}
	if m.A != 1 || m.D != 1 || m.B != 0 || m.C != 0 {
		return false
	}
	x, y := -int(math.Floor(0.5-m.E)), -int(math.Floor(0.5-m.F))
	if !Hides(Placed{s.Image, x, y}, g.depth(o), g.covers) {
		return false
	}
	BlendFaded(cv, Placed{s.Image, x, y})
	return true
}
