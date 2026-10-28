package tilemap

import "github.com/gabrielluizsf/sge/collider"

// Solid is what Move walks through: a map, or anything else with walls,
// such as terrain made up as it is explored.
type Solid interface {
	// Blocks reports whether c, a collider in the world, touches a wall.
	Blocks(c collider.Collider) bool
	// Stride is the most a mover may go in one step without passing
	// through a wall: half the thinnest wall.
	Stride() float64
}
