package farm

import (
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// seed is the wild's: the same seed, the same land around the farm.
const seed = 20261024

// world is the land the farm lies in: its tiles, the endless wild, the
// towns in it, and the art it is drawn in.
type world struct {
	art   Art
	tiles *tilemap.TileSet
	wild  *terrain.World
	towns *town.Layout
}

// newWild is the endless land around the farm, as areas over grass: lakes
// with shores, and trees, thick in woods and now and then alone, kept off
// the water's edge; the towns, their streets and gardens, their houses'
// cells and the smithy's, and the road to the farm, kept clear of both.
// Nothing of it lies on the farm, nor beside its edge.
func newWild(w world, farmSize [2]int) *terrain.World {
	tiles, towns := w.tiles, w.towns
	cols, rows := farmSize[0], farmSize[1]
	onFarm := func(c, r int) bool { return c >= 0 && r >= 0 && c < cols && r < rows }
	edge := 2 // a look with trees grown large keeps them further off the farm's fence
	if w.art.Reach() > 0 {
		edge = 16
	}
	nearFarm := func(c, r int) bool { return c >= -edge && r >= -edge && c < cols+edge && r < rows+edge }
	// The roads between the towns keep off the farm and its edge.
	towns.Keep = []town.Rect{{C0: -edge, R0: -edge, C1: cols + edge, R1: rows + edge}}
	built := func(c, r, margin int) bool { return towns.Grounds(c, r, margin) || towns.OnRoad(c, r, margin) }
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
		return towns.OnRoad(c, r, 0) || t != nil && t.Road(c, r)
	}, Tile: func(int, int) int { return tileID["path"] }})
	gen.Add(terrain.Area{Name: "fences", Layer: wildThings, In: towns.Fence,
		Tile: func(c, r int) int {
			if towns.FenceDown(c, r) {
				return tileID[fenceDown]
			}
			return tileID["fence"]
		}})
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
	// Trees grown to their real size need room: a look whose pictures are
	// larger than a cell spaces them wider.
	thickGap, loneGap, keepOff := 1, 2, 1
	if w.art.Reach() > 0 {
		thickGap, loneGap, keepOff = 3, 4, 6
	}
	thick := terrain.Scatter(seed+3, 0.5, thickGap)
	lone := terrain.Scatter(seed+4, 0.04, loneGap)
	tree := func(c, r int) bool {
		if !(woods(c, r) && thick(c, r) || lone(c, r)) {
			return false
		}
		return !nearFarm(c, r) && !built(c, r, keepOff) && !gen.Near("water", c, r, 1)
	}
	gen.Add(terrain.Area{Name: "trees", Layer: wildThings, In: tree, Tile: treeAt})
	// Rocks lie by biome, apart from each other, the trees and the water,
	// off the farm, the towns and the roads.
	rocks := terrain.Scatter(seed+6, rockCeiling, rockGap(w.art))
	rock := func(c, r int) bool {
		if !rocks(c, r) || terrain.Chance(seed+13, c, r)*rockCeiling >= rockSpace(c, r) {
			return false
		}
		return !nearFarm(c, r) && !built(c, r, 2) && !gen.Near("water", c, r, 1) && !tree(c, r)
	}
	gen.Add(terrain.Area{Name: "rocks", Layer: wildThings, In: rock, Tile: rockAt})
	gen.Add(terrain.Area{Name: "buildings", Layer: wildThings, In: func(c, r int) bool { return buildingAt(w, c, r) >= 0 },
		Tile: func(c, r int) int { return buildingAt(w, c, r) }})
	wild := terrain.New("wild", gen, 16, zGround, zWater, zThings)
	// Its trees stand among the farmer and each other.
	wild.Sorted = []int{wildThings}
	return wild
}

// buildingAt is the tile of what is built on cell col, row of a town: the
// cells of a house or of the stall, solid and drawn by their objects, the
// anvil or the forge; -1 where nothing is.
func buildingAt(w world, col, row int) int {
	t := w.towns.At(col, row)
	if t == nil {
		return -1
	}
	c := town.Cell{Col: col, Row: row}
	switch {
	case c == t.Anvil:
		return tileID["anvil"]
	case c == t.Forge:
		return tileID["forge"]
	case stallRect(w.art, t.Stall).has(col, row) || wellRect(w.art, t).has(col, row) || t.House(col, row) != nil:
		return tileID["block"]
	case stableRect(w.art, t).has(col, row):
		return tileID["block"]
	}
	return fenceAt(townPen(t), col, row)
}
