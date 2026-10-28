package iso

import "github.com/gabrielluizsf/sge/vector"

// Project is where the world point p falls on the screen, from where the
// world's origin does. A point's Z lifts it up the screen.
func (p Projection) Project(w vector.Vector3) vector.Vector3 {
	hx, hy := p.Width/2/p.Cell, p.Height/2/p.Cell
	return vector.Vector3{X: (w.X - w.Y) * hx, Y: (w.X+w.Y)*hy - w.Z*2*hy}
}
