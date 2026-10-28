// Package pathfind finds ways across a grid of cells, around what stands
// in the way, and steers along them: villagers walking home, enemies
// coming round a wall.
package pathfind

// Cell is a cell of the grid.
type Cell struct{ Col, Row int }

// Grid is the cells one may walk: Open reports whether a cell may be
// stood on, and Limit is how many cells a search looks at before giving
// up, so a way to somewhere walled off costs no more than that.
type Grid struct {
	Open  func(col, row int) bool
	Limit int
}

// reach is how far from the origin, in cells, a grid's cells are keyed.
const reach = 1 << 20

// key is c as one number, which JavaScript maps compare fast; ok is false
// beyond reach.
func key(c Cell) (k int, ok bool) {
	if c.Col <= -reach || c.Col >= reach || c.Row <= -reach || c.Row >= reach {
		return 0, false
	}
	return (c.Col+reach)*(2*reach) + c.Row + reach, true
}
