package terrain

// chunkSize is a chunk's width and height in the world.
func (w *World) chunkSize() (float64, float64) {
	return float64(w.Size * w.Gen.Tiles.Width), float64(w.Size * w.Gen.Tiles.Height)
}
