package main

import (
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// seed is the wild's: the same seed, the same land around the farm.
const seed = 20261024

// world is the land the farm lies in: its tiles, the endless wild, and the
// towns in it.
type world struct {
	tiles *tilemap.TileSet
	wild  *terrain.World
	towns *town.Layout
}

// newWild is the endless land around the farm, as areas over grass: lakes
// with shores, and trees, thick in woods and now and then alone, kept off
// the water's edge; the towns, their streets and gardens, their houses'
// cells and the smithy's, and the road to the farm, kept clear of both.
// Nothing of it lies on the farm, nor beside its edge.
func newWild(tiles *tilemap.TileSet, cols, rows int, towns *town.Layout) *terrain.World {
	onFarm := func(c, r int) bool { return c >= 0 && r >= 0 && c < cols && r < rows }
	nearFarm := func(c, r int) bool { return c >= -2 && r >= -2 && c < cols+2 && r < rows+2 }
	built := func(c, r, margin int) bool { return towns.Grounds(c, r, margin) || highway(c, r, margin) }
	gen := &terrain.Generator{Tiles: tiles, Layers: 3}
	gen.Add(terrain.Area{Name: "grass", Layer: wildGround,
		In: func(c, r int) bool { return !onFarm(c, r) },
		Tile: func(c, r int) int {
			if terrain.Chance(seed, c, r) < 0.06 {
				return tileID["flowers0"]
			}
			return tileID["grass0"]
		}})
	gen.Add(terrain.Area{Name: "road", Layer: wildGround, In: func(c, r int) bool {
		t := towns.At(c, r)
		return highway(c, r, 0) || t != nil && t.Road(c, r)
	}, Tile: func(int, int) int { return tileID["path"] }})
	gen.Add(terrain.Area{Name: "gardens", Layer: wildGround, In: func(c, r int) bool {
		t := towns.At(c, r)
		if t == nil {
			return false
		}
		for _, h := range t.Houses {
			if h.Garden.Col == c && h.Garden.Row == r {
				return true
			}
		}
		return false
	}, Tile: func(int, int) int { return tileID["soil_dry"] }})
	// Lakes are kept off the farm and the towns before they are opened, so
	// where one is cut short no thread of it is left. Opening asks about
	// each cell many times, so the answers are kept while a chunk is made.
	lakes := terrain.Blob(seed+1, 9, 0.6)
	raw := terrain.Memo(func(c, r int) bool { return lakes(c, r) && !nearFarm(c, r) && !built(c, r, 2) }, 1<<14)
	water := terrain.Memo(terrain.Open(raw), 1<<14)
	gen.Add(terrain.Area{Name: "water", Layer: wildWater, In: water, Tile: terrain.Edged(water, shore)})
	woods := terrain.Blob(seed+2, 12, 0.55)
	thick := terrain.Scatter(seed+3, 0.5, 1)
	lone := terrain.Scatter(seed+4, 0.04, 2)
	tree := func(c, r int) bool {
		if !(woods(c, r) && thick(c, r) || lone(c, r)) {
			return false
		}
		return !nearFarm(c, r) && !built(c, r, 1) && !gen.Near("water", c, r, 1)
	}
	gen.Add(terrain.Area{Name: "trees", Layer: wildThings, In: tree, Tile: func(int, int) int { return tileID["tree0"] }})
	gen.Add(terrain.Area{Name: "buildings", Layer: wildThings, In: func(c, r int) bool { return buildingAt(towns, c, r) >= 0 },
		Tile: func(c, r int) int { return buildingAt(towns, c, r) }})
	w := terrain.New("wild", gen, 16, zGround, zWater, zThings)
	// Its trees stand among the farmer and each other.
	w.Sorted = []int{wildThings}
	return w
}

// buildingAt is the tile of what is built on cell col, row of a town: a
// house's or the stall's cells, solid and drawn by their objects, the
// anvil or the forge; -1 where nothing is.
func buildingAt(towns *town.Layout, col, row int) int {
	t := towns.At(col, row)
	if t == nil {
		return -1
	}
	c := town.Cell{Col: col, Row: row}
	switch {
	case c == t.Anvil:
		return tileID["anvil"]
	case c == t.Forge:
		return tileID["forge"]
	case c == t.Stall || t.House(col, row) != nil:
		return tileID["block"]
	}
	return -1
}
