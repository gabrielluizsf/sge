package terrain

// chunkSize is a chunk's width and height in the world.
func (w *World) chunkSize() (float64, float64) {
	cw, ch := w.Gen.Tiles.Cell()
	return float64(w.Size * cw), float64(w.Size * ch)
}
