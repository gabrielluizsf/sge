package terrain

import "math"

// CellOf is the cell under the world point x, y.
func (w *World) CellOf(x, y float64) (col, row int) {
	cw, ch := w.Gen.Tiles.Cell()
	return int(math.Floor(x / float64(cw))), int(math.Floor(y / float64(ch)))
}
