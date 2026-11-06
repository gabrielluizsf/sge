package terrain

// warmUp gives the warm queue the pictures of tile id, the first time a
// chunk has it, if the world has a queue.
func (w *World) warmUp(id int, seen map[int]bool) {
	if w.Warm == nil || seen[id] {
		return
	}
	seen[id] = true
	w.Warm.Add(w.Gen.Tiles.PicturesOf(id)...)
}
