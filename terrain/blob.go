package terrain

// Blob is the cells where Fractal noise of the given scale, in cells, is
// above level: round patches some scale across, as lakes or woods are,
// fewer and smaller the higher level is.
func Blob(seed uint32, scale, level float64) func(col, row int) bool {
	return func(col, row int) bool {
		return Fractal(seed, float64(col)/scale, float64(row)/scale, 3) > level
	}
}
