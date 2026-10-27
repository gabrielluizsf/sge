package terrain

import "math"

// InAt reports whether the world point x, y lies in the area called name:
// whether one standing there is in the water, say. It asks the generator,
// so it holds for any place, made into a chunk yet or not.
func (w *World) InAt(name string, x, y float64) bool {
	cw, ch := w.Gen.Tiles.Cell()
	return w.Gen.In(name, int(math.Floor(x/float64(cw))), int(math.Floor(y/float64(ch))))
}
