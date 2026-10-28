package terrain

// TileAt is the tile cell col, row shows on layer: read from its chunk
// when that is made, which is fast, and worked out by the generator
// otherwise, edits and all.
func (w *World) TileAt(layer, col, row int) int {
	cx, cy := floorDiv(col, w.Size), floorDiv(row, w.Size)
	if key, ok := cellKey(cx, cy); ok {
		if c := w.chunks[key]; c != nil && layer >= 0 && layer < len(c.layers) {
			return c.layers[layer].At(col-cx*w.Size, row-cy*w.Size)
		}
	}
	return w.Gen.Tile(layer, col, row)
}

// floorDiv is a divided by b, rounded down, negative a and all.
func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && (a < 0) != (b < 0) {
		q--
	}
	return q
}
