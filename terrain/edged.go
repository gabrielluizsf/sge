package terrain

// Edged is the Tile of an area whose tiles join up, a lake's shores: the
// tile for each cell is tile of its Mask.
func Edged(in func(col, row int) bool, tile func(mask int) int) func(col, row int) int {
	return func(col, row int) int {
		return tile(Mask(in, col, row))
	}
}
