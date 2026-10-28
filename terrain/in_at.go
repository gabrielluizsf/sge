package terrain

// InAt reports whether the world point x, y lies in the area called name:
// whether one standing there is in the water, say. It asks the generator,
// so it holds for any place, made into a chunk yet or not.
func (w *World) InAt(name string, x, y float64) bool {
	col, row := w.CellOf(x, y)
	return w.Gen.In(name, col, row)
}
