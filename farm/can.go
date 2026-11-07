package farm

import "github.com/gabrielluizsf/sge/tilemap"

// canCapacity is how much water the can holds: a watering takes one.
const canCapacity = 20

// inCan is how much water is in k's can.
func (k *kit) inCan() int { return canCapacity - k.used }

// pour uses a watering of k's can's water, and reports whether it had
// any.
func (k *kit) pour() bool {
	if k.used >= canCapacity {
		return false
	}
	k.used++
	return true
}

// fill fills k's can.
func (k *kit) fill() { k.used = 0 }

// isWater reports whether cell col, row is water, a pond's or a shore on
// the farm, a lake's in the wild, or a town's well: where the can is
// filled.
func (f *farm) isWater(col, row int) bool {
	if f.onFarm(col, row) {
		return f.water.At(col, row) != tilemap.Empty
	}
	if t := f.towns.At(col, row); t != nil && atWell(wellRect(f.art, t), col, row) {
		return true
	}
	return f.wild.Gen.In("water", col, row)
}

// atWell reports whether cell col, row is the well's of w, which covers
// at least its own corner cell, however small the art draws it.
func atWell(w rect, col, row int) bool {
	w.cols, w.rows = max(w.cols, 1), max(w.rows, 1)
	return w.has(col, row)
}
