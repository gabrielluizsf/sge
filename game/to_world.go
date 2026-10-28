package game

import "github.com/gabrielluizsf/sge/vector"

// ToWorld is the point of the ground drawn at the screen's point s: what
// the mouse is over.
func (g *Game) ToWorld(s vector.Vector3) vector.Vector3 {
	return g.unproject(s.Add(g.camera()))
}
