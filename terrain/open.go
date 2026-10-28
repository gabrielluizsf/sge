package terrain

// Open is in with everything thinner than three cells taken away: a cell
// is in it when some three by three square of cells around it lies wholly
// in in. Lakes and woods so opened have no threads, no lone cells and no
// nubs, whatever the noise under them.
func Open(in func(col, row int) bool) func(col, row int) bool {
	whole := func(col, row int) bool {
		for r := row - 1; r <= row+1; r++ {
			for c := col - 1; c <= col+1; c++ {
				if !in(c, r) {
					return false
				}
			}
		}
		return true
	}
	return func(col, row int) bool {
		for r := row - 1; r <= row+1; r++ {
			for c := col - 1; c <= col+1; c++ {
				if whole(c, r) {
					return true
				}
			}
		}
		return false
	}
}
