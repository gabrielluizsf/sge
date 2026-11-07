package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/farm/art"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// TileNames are the tiles the maps are made of, each of which an Art
// draws, in the order of their ids in the tile set: the ground, the
// farm's things, the crops, what stands in the towns, the houses'
// insides, and the shores of water, every mask four frames into its
// ripple.
func TileNames() []string {
	names := []string{
		"grass0", "grass1", "grass2", "grass3",
		"flowers0", "flowers1", "flowers2", "flowers3",
		"path", "soil_dry", "soil_wet", "fence", "block",
		"stump", "anvil", witheredTile,
	}
	for _, c := range crops {
		for stage := range 5 {
			names = append(names, fmt.Sprint(c.name, stage))
		}
	}
	for _, biome := range TreeBiomes {
		for v := range TreeVariants {
			for f := range TreeFrames {
				names = append(names, treeName(biome, v, f))
			}
		}
	}
	for _, biome := range TreeBiomes {
		for v := range RockVariants {
			names = append(names, rockName(biome, v))
		}
	}
	names = append(names, interiorNames...)
	for _, mask := range terrain.Masks() {
		for f := range 4 {
			names = append(names, shoreName(mask, f))
		}
	}
	// Last, so the ids of the trees and rocks a save names stay as they were.
	names = append(names, fenceDown)
	return append(names, fenceCorners[:]...)
}

// fenceDown is the tile of a fence turned to run down the rows.
const fenceDown = "fence_down"

// fenceCorners are the tiles of a fence's corners, by the corner of a pen
// each is: its rails run from the post in the middle of its cell to the next
// cells of the fence on either side.
var fenceCorners = [4]string{"fence_nw", "fence_ne", "fence_sw", "fence_se"}

// shoreName is the sprite of water f steps into its ripple, with water
// beside it as mask says.
func shoreName(mask, f int) string { return fmt.Sprint("shore", mask, "_", f) }

// shore is the tile of water with water beside it as mask says.
func shore(mask int) int { return tileID[shoreName(mask, 0)] }

// isFence reports whether id is a fence's tile, turned either way.
func isFence(id int) bool {
	if id == tileID["fence"] || id == tileID[fenceDown] {
		return true
	}
	for _, c := range fenceCorners {
		if id == tileID[c] {
			return true
		}
	}
	return false
}

// tileID is the id of each tile in the set, by name.
var tileID = map[string]int{}

// newTiles is the tile set, each tile drawn by a the first time a map
// shows it, and which are solid and which sway. Each tile's picture is
// a cell's diamond, isoW by tile, at the bottom of isoH, what stands on
// the cell reaching up out of it; its cell in the world stays tile square.
func newTiles(a Art) *tilemap.TileSet {
	names := append(TileNames(), pieceTileNames()...)
	for i, name := range names {
		tileID[name] = i
	}
	tiles := tilemap.NewLazyTileSet(tilemap.Lazy{Count: len(names), Width: isoW, Height: isoH,
		Make: func(id int) *canvas.Canvas { return drawTile(a, names[id]) }})
	tiles.CellWidth, tiles.CellHeight = tile, tile
	tiles.Reach = a.Reach()
	id := func(names ...string) []int {
		ids := make([]int, len(names))
		for i, n := range names {
			ids[i] = tileID[n]
		}
		return ids
	}
	tiles.SetSolid(true, id("fence", fenceDown, "block", "stump", "anvil")...)
	tiles.SetSolid(true, id(fenceCorners[:]...)...)
	tiles.SetSolid(true, id(interiorNames[len(interiorFloors):]...)...)
	tiles.SetSolid(true, id(pieceTileNames()...)...)
	// Grass, flowers and trees sway each in their own time; water ripples
	// all together, shores and all.
	gridded := a.Gridded()
	tiles.Animate(tileID["grass0"], tilemap.Animation{Frames: id("grass0", "grass1", "grass2", "grass3"), FPS: 2, Staggered: true, Parity: gridded})
	tiles.Animate(tileID["flowers0"], tilemap.Animation{Frames: id("flowers0", "flowers1", "flowers2", "flowers3"), FPS: 2, Staggered: true, Parity: gridded})
	for _, biome := range TreeBiomes {
		for v := range TreeVariants {
			frames := id(treeName(biome, v, 0), treeName(biome, v, 1), treeName(biome, v, 2), treeName(biome, v, 3))
			treeIDs[frames[0]] = true
			tiles.SetSolid(true, frames[0])
			tiles.Animate(frames[0], tilemap.Animation{Frames: frames, FPS: 1.5, Staggered: true})
		}
	}
	for _, biome := range TreeBiomes {
		for v := range RockVariants {
			id := tileID[rockName(biome, v)]
			rockIDs[id] = true
			tiles.SetSolid(true, id)
		}
	}
	// Water is no wall: the farmer swims in it.
	for _, mask := range terrain.Masks() {
		tiles.Animate(shore(mask), tilemap.Animation{Frames: id(shoreName(mask, 0), shoreName(mask, 1), shoreName(mask, 2), shoreName(mask, 3)), FPS: 4})
	}
	// A ripe crop sways, asking to be picked.
	for _, c := range crops {
		ripe3, ripe4 := fmt.Sprint(c.name, 3), fmt.Sprint(c.name, 4)
		tiles.Animate(tileID[ripe3], tilemap.Animation{Frames: id(ripe3, ripe4), FPS: 2, Staggered: true})
	}
	return tiles
}

// interiorFloors are the floors of the houses' insides, and interiorNames
// all their tiles, the floors first: what stands on them is solid.
var (
	interiorFloors = []string{"floor", "mat", "rug0", "rug1", "rug2", "rug3"}
	interiorNames  = append(append(append([]string{}, interiorFloors...),
		"forge", "bed", "bed_blue", "bed_double", "table", "chair", "shelf", "stove", "plant", "chest"), wallNames()...)
)

// wallNames are the back walls of each kind of house's inside: along a
// cell's north edge, the same with a window, along its west, and both.
func wallNames() []string {
	var out []string
	for k := range art.WallKinds {
		out = append(out, fmt.Sprint("wall_n", k), fmt.Sprint("wall_nwin", k), fmt.Sprint("wall_w", k), fmt.Sprint("wall_nw", k))
	}
	return out
}
