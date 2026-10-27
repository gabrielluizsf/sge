package game

import "github.com/gabrielluizsf/sge/vector"

// ToScreen is where the world point p is drawn, as the camera sees it: a
// painter of an object in the world draws there.
func (g *Game) ToScreen(p vector.Vector3) vector.Vector3 {
	return g.project(p).Sub(g.camera())
}
