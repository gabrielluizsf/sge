package iso

import "github.com/gabrielluizsf/sge/vector"

// Unproject is the point of the ground, Z 0, that falls at s on the
// screen, from where the world's origin does: the point under a mouse.
func (p Projection) Unproject(s vector.Vector3) vector.Vector3 {
	a, b := s.X/(p.Width/2/p.Cell), s.Y/(p.Height/2/p.Cell)
	return vector.Vector3{X: (a + b) / 2, Y: (b - a) / 2}
}
