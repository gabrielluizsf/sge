package terrain

// Scatter is cells strewn with the given chance each, none within spacing
// cells of another: trees in a wood, say, never trunk to trunk. Of the
// cells that would be too close, the one with the smallest hash is kept,
// which any cell can work out alone, so chunks made apart agree.
func Scatter(seed uint32, chance float64, spacing int) func(col, row int) bool {
	limit := uint32(chance * (1 << 32))
	if chance >= 1 {
		limit = ^uint32(0)
	}
	return func(col, row int) bool {
		h := Hash(seed, col, row)
		if h >= limit {
			return false
		}
		for r := row - spacing; r <= row+spacing; r++ {
			for c := col - spacing; c <= col+spacing; c++ {
				if o := Hash(seed, c, r); o < h || o == h && (r < row || r == row && c < col) {
					return false
				}
			}
		}
		return true
	}
}
