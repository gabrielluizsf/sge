package farm

import (
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// swimSpeed is how fast the farmer swims: slower than they walk. They
// wade where the water is shallow, at wadeSpeed, between the two.
const (
	swimSpeed = walkSpeed * 0.55
	wadeSpeed = walkSpeed * 0.75
)

// deepFrom is how many cells from any land water is deep enough to swim
// in: the middle of a river or a lake. Nearer the land it is shallow, and
// one wades.
const deepFrom = 2

// inWater reports whether the point p, where one stands, is in water: the
// farm's pond, or a lake of the wild.
func (f *farm) inWater(p vector.Vector3) bool {
	if f.room != nil {
		return false
	}
	return f.water.TileAt(p.X, p.Y) != tilemap.Empty || f.wild.InAt("water", p.X, p.Y)
}

// wet reports whether cell col, row is water: the farm's pond, or a lake
// of the wild.
func (f *farm) wet(col, row int) bool {
	if f.onFarm(col, row) {
		return f.water.At(col, row) != tilemap.Empty
	}
	return f.wild.Gen.In("water", col, row)
}

// deepAt reports whether the point p, where one stands, is in deep water:
// a water cell with no land within deepFrom-1 cells of it, nowhere near
// the shore. Only there does one swim.
func (f *farm) deepAt(p vector.Vector3) bool {
	if f.room != nil {
		return false
	}
	col, row := f.cellOf(p)
	for dc := -(deepFrom - 1); dc < deepFrom; dc++ {
		for dr := -(deepFrom - 1); dr < deepFrom; dr++ {
			if !f.wet(col+dc, row+dr) {
				return false
			}
		}
	}
	return true
}

// solids are what one walks into where the player is: a room's walls and
// furniture, or the farm's and the wild's.
func (f *farm) solids() []tilemap.Solid {
	if f.room != nil {
		return []tilemap.Solid{f.room.things}
	}
	return []tilemap.Solid{f.water, f.things, f.wild, chestWalls{f}}
}
