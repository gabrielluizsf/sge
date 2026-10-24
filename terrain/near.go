package terrain

// Near reports whether a cell within dist cells of col, row is in the area
// called name.
func (g *Generator) Near(name string, col, row, dist int) bool {
	for r := row - dist; r <= row+dist; r++ {
		for c := col - dist; c <= col+dist; c++ {
			if g.In(name, c, r) {
				return true
			}
		}
	}
	return false
}
