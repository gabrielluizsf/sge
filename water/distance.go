package water

import (
	"math"

	"github.com/gabrielluizsf/sge/terrain"
)

// ring are the neighbours of a cell, each mask bit's with its step from
// the cell: columns along u, rows along v.
var ring = [8]struct{ bit, du, dv int }{
	{terrain.North, 0, -1}, {terrain.East, 1, 0}, {terrain.South, 0, 1}, {terrain.West, -1, 0},
	{terrain.NorthEast, 1, -1}, {terrain.SouthEast, 1, 1}, {terrain.SouthWest, -1, 1}, {terrain.NorthWest, -1, -1},
}

// landOf are the neighbours of a cell that are land, by mask: the ones
// whose bit is not set.
func landOf(mask int) [][2]float64 {
	var out [][2]float64
	for _, r := range ring {
		if mask&r.bit == 0 {
			out = append(out, [2]float64{float64(r.du), float64(r.dv)})
		}
	}
	return out
}

// distance is how far (u, v) is from the nearest of the land cells, in
// cells; each is a unit square at its step from the cell. Land two cells
// away is not counted, which is far beyond where the shore matters.
func distance(land [][2]float64, u, v float64) float64 {
	best := 1e9
	for _, c := range land {
		dx := math.Max(math.Max(c[0]-u, u-c[0]-1), 0)
		dy := math.Max(math.Max(c[1]-v, v-c[1]-1), 0)
		best = math.Min(best, math.Hypot(dx, dy))
	}
	return best
}
