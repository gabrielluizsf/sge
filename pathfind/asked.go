package pathfind

// asked is the grid's Open that remembers its answers for one search: a
// cell is asked about from each of its neighbours, and working out whether
// it may be stood on can cost much more than looking it up.
func (g Grid) asked() func(col, row int) bool {
	seen := map[int]bool{}
	return func(col, row int) bool {
		k, ok := key(Cell{Col: col, Row: row})
		if !ok {
			return g.Open(col, row)
		}
		open, ok := seen[k]
		if !ok {
			open = g.Open(col, row)
			seen[k] = open
		}
		return open
	}
}
