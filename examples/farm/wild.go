package main

import (
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// seed is the wild's: the same seed, the same land around the farm.
const seed = 20261024

// newWild is the endless land around the farm, as areas over grass: lakes
// with shores, and trees, thick in woods and now and then alone, kept off
// the water's edge. Nothing of it lies on the farm, nor beside its edge.
func newWild(tiles *tilemap.TileSet, cols, rows int) *terrain.World {
	onFarm := func(c, r int) bool { return c >= 0 && r >= 0 && c < cols && r < rows }
	nearFarm := func(c, r int) bool { return c >= -2 && r >= -2 && c < cols+2 && r < rows+2 }
	gen := &terrain.Generator{Tiles: tiles, Layers: 4}
	gen.Add(terrain.Area{Name: "grass", Layer: 0,
		In: func(c, r int) bool { return !onFarm(c, r) },
		Tile: func(c, r int) int {
			if terrain.Chance(seed, c, r) < 0.06 {
				return tileID["flowers0"]
			}
			return tileID["grass0"]
		}})
	// Lakes are kept off the farm before they are opened, so where one is
	// cut short by the farm no thread of it is left. Opening asks about
	// each cell many times, so the answers are kept while a chunk is made.
	lakes := terrain.Blob(seed+1, 9, 0.6)
	raw := terrain.Memo(func(c, r int) bool { return lakes(c, r) && !nearFarm(c, r) }, 1<<14)
	water := terrain.Memo(terrain.Open(raw), 1<<14)
	gen.Add(terrain.Area{Name: "water", Layer: 1, In: water, Tile: terrain.Edged(water, shore)})
	woods := terrain.Blob(seed+2, 12, 0.55)
	thick := terrain.Scatter(seed+3, 0.5, 1)
	lone := terrain.Scatter(seed+4, 0.04, 2)
	tree := func(c, r int) bool {
		if !(woods(c, r) && thick(c, r) || lone(c, r)) {
			return false
		}
		// Its crown, the cell above, must be off the farm too.
		return !nearFarm(c, r) && !nearFarm(c, r-1) && !gen.Near("water", c, r, 1)
	}
	gen.Add(terrain.Area{Name: "trees", Layer: 2, In: tree, Tile: func(int, int) int { return tileID["tree_trunk"] }})
	gen.Add(terrain.Area{Name: "crowns", Layer: 3, In: func(c, r int) bool { return tree(c, r+1) },
		Tile: func(int, int) int { return tileID["tree_crown0"] }})
	return terrain.New("wild", gen, 16, zGround, zWater, zThings, zOver)
}
