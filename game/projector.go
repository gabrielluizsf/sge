package game

import "github.com/gabrielluizsf/sge/vector"

// A Projector turns the world into what the screen shows: an isometric
// view (package iso), say, of a world that is a flat grid underneath.
// Both are linear: Project maps a step in the world to one on the screen,
// Unproject a step on the screen to one on the ground.
type Projector interface {
	Project(world vector.Vector3) vector.Vector3
	Unproject(screen vector.Vector3) vector.Vector3
}

// project is p as the game's projection draws it.
func (g *Game) project(p vector.Vector3) vector.Vector3 {
	if g.Projection == nil {
		return p
	}
	return g.Projection.Project(p)
}

// unproject is the step on the ground the screen's step s shows.
func (g *Game) unproject(s vector.Vector3) vector.Vector3 {
	if g.Projection == nil {
		return s
	}
	return g.Projection.Unproject(s)
}
