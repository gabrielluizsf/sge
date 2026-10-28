package terrain

// The neighbours of a cell, as the bits of a mask: a bit is set when the
// neighbour that way is of the same area. Tiles that join up, a lake's
// shores, are chosen by mask.
const (
	North = 1 << iota
	East
	South
	West
	NorthEast
	SouthEast
	SouthWest
	NorthWest
	// AllSides is a cell inside its area, with it all round.
	AllSides = North | East | South | West | NorthEast | SouthEast | SouthWest | NorthWest
)

// neighbours are each bit's cell, as a step from the cell.
var neighbours = [8][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}

// Mask is which neighbours of cell col, row are in the area in reports,
// Canonical.
func Mask(in func(col, row int) bool, col, row int) int {
	mask := 0
	for i, d := range neighbours {
		if in(col+d[0], row+d[1]) {
			mask |= 1 << i
		}
	}
	return Canonical(mask)
}
