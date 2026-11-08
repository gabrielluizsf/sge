package terrain

import "math"

// Biome is a kind of land, a forest or a coast, and the climate it is
// found in: a warmth and a wetness, each from 0 to 1.
type Biome struct {
	Name            string
	Warmth, Wetness float64
}

// Climate lays biomes over an endless world, as ecologists chart them by
// climate: a warmth and a wetness that wander across it as Fractal noise
// of Scale cells, from Seed, each cell taking the biome whose climate is
// nearest its own. Ragged is how far, as climate, a cell's own is shaken,
// so biomes meet in a ragged fringe, not along a line.
type Climate struct {
	Seed   uint32
	Scale  float64
	Ragged float64
	Biomes []Biome
}

// climate is cell col, row's warmth and wetness, spread from the noise's
// crowded middle across 0 to 1, and shaken.
func (c Climate) climate(col, row int) (warm, wet float64) {
	x, y := float64(col)/c.Scale, float64(row)/c.Scale
	spread := func(v float64) float64 { return math.Max(0, math.Min(1, (v-0.5)*2.4+0.5)) }
	warm = spread(Fractal(c.Seed, x, y, 3)) + (Chance(c.Seed+11, col, row)-0.5)*c.Ragged
	wet = spread(Fractal(c.Seed+7, x, y, 3)) + (Chance(c.Seed+13, col, row)-0.5)*c.Ragged
	return warm, wet
}

// At is the index, in Biomes, of cell col, row's biome; -1 with none.
func (c Climate) At(col, row int) int {
	warm, wet := c.climate(col, row)
	best, far := -1, math.Inf(1)
	for i, b := range c.Biomes {
		dw, dm := warm-b.Warmth, wet-b.Wetness
		if d := dw*dw + dm*dm; d < far {
			best, far = i, d
		}
	}
	return best
}

// In is the cells of the biome called name, for an Area's In.
func (c Climate) In(name string) func(col, row int) bool {
	want := -1
	for i, b := range c.Biomes {
		if b.Name == name {
			want = i
		}
	}
	return func(col, row int) bool { return want >= 0 && c.At(col, row) == want }
}
