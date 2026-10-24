package game

import (
	"math"

	"github.com/gabrielluizsf/sge/vector"
)

// camera is where the world is seen from, to whole pixels: everything
// moves by the same amount, so what meets in the world meets on screen.
func (g *Game) camera() vector.Vector3 {
	return vector.Vector3{X: math.Floor(g.Camera.X + 0.5), Y: math.Floor(g.Camera.Y + 0.5)}
}
