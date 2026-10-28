package terrain

// Set edits cell col, row of layer to tile id, as Generator.Edit does,
// and shows it at once if its chunk is made.
func (w *World) Set(layer, col, row, id int) {
	w.Gen.Edit(layer, col, row, id)
	cx, cy := floorDiv(col, w.Size), floorDiv(row, w.Size)
	key, ok := cellKey(cx, cy)
	if !ok {
		return
	}
	if c := w.chunks[key]; c != nil && layer >= 0 && layer < len(c.layers) {
		c.layers[layer].Set(col-cx*w.Size, row-cy*w.Size, id)
	}
}
