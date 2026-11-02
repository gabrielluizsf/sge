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
// the farm or a lake's in the wild: where the can is filled.
func (f *farm) isWater(col, row int) bool {
	if f.onFarm(col, row) {
		return f.water.At(col, row) != tilemap.Empty
	}
	return f.wild.Gen.In("water", col, row)
}
