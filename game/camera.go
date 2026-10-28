package game

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// camera is where the camera's point of the world falls on the screen, to
// whole pixels: everything moves by the same amount, so what meets in the
// world meets on screen.
func (g *Game) camera() vector.Vector3 {
	c := g.project(g.Camera)
	return vector.Vector3{X: math.Floor(c.X + 0.5), Y: math.Floor(c.Y + 0.5)}
}
