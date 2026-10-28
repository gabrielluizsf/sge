package pathfind

import "github.com/gabrielluizsf/sge/vector"

// Route is a path being walked: its cells, and how many have been
// reached.
type Route struct {
	Cells   []Cell
	reached int
}

// arrived is how near, in the world's units, a cell's point is reached.
const arrived = 0.5

// Steer is the way from at to the route's next cell's point, where
// center puts it, as long as it is far: what is left to walk to it. Cells
// reached are passed, so the walk goes on without a stop; done is true
// when the last one is reached.
func (r *Route) Steer(at vector.Vector3, center func(Cell) vector.Vector3) (way vector.Vector3, done bool) {
	for r.reached < len(r.Cells) {
		way = center(r.Cells[r.reached]).Sub(at)
		way.Z = 0
		if way.Length() > arrived {
			return way, false
		}
		r.reached++
	}
	return vector.Vector3{}, true
}
