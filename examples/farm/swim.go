package main

import (
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// swimSpeed is how fast the farmer swims: slower than they walk.
const swimSpeed = walkSpeed * 0.55

// inWater reports whether the point p, where one stands, is in water: the
// farm's pond, or a lake of the wild.
func (f *farm) inWater(p vector.Vector3) bool {
	return f.water.TileAt(p.X, p.Y) != tilemap.Empty || f.wild.InAt("water", p.X, p.Y)
}
