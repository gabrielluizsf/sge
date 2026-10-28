package terrain

// Smooth is in with its ragged bits worn away: a cell is in it when at
// least need of the nine cells around it, itself among them, are in in.
// With need 5, a majority, threads a cell wide and lone cells go and
// notches fill, so lakes and woods have round, solid shapes.
func Smooth(in func(col, row int) bool, need int) func(col, row int) bool {
	return func(col, row int) bool {
		n := 0
		for r := row - 1; r <= row+1; r++ {
			for c := col - 1; c <= col+1; c++ {
				if in(c, r) {
					if n++; n >= need {
						return true
					}
				}
			}
		}
		return false
	}
}
