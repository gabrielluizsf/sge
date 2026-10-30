package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// tileNames are the sprites the maps are made of, in the order of their
// ids in the tile set: the ground, the farm's things, the crops, what
// stands in the towns, the houses' insides, and the shores of water, every mask four frames into its ripple.
func tileNames() []string {
	names := []string{
		"grass0", "grass1", "grass2", "grass3",
		"flowers0", "flowers1", "flowers2", "flowers3",
		"path", "soil_dry", "soil_wet", "fence", "block",
		"tree0", "tree1", "tree2", "tree3",
		"carrot0", "carrot1", "carrot2", "carrot3", "carrot4",
		"tomato0", "tomato1", "tomato2", "tomato3", "tomato4",
		"stump", "anvil",
	}
	names = append(names, interiorNames...)
	for _, mask := range terrain.Masks() {
		for f := range 4 {
			names = append(names, shoreName(mask, f))
		}
	}
	return names
}

// shoreName is the sprite of water f steps into its ripple, with water
// beside it as mask says.
func shoreName(mask, f int) string { return fmt.Sprint("shore", mask, "_", f) }

// shore is the tile of water with water beside it as mask says.
func shore(mask int) int { return tileID[shoreName(mask, 0)] }

// tileID is the id of each tile in the set, by name.
var tileID = map[string]int{}

// newTiles is the tile set, each tile drawn from its SVG the first time a
// map shows it, and which are solid and which sway. Each tile's picture is
// a cell's diamond, isoW by tile, at the bottom of isoH, what stands on
// the cell reaching up out of it; its cell in the world stays tile square.
func newTiles() *tilemap.TileSet {
	names := tileNames()
	for i, name := range names {
		tileID[name] = i
	}
	tiles := tilemap.NewLazyTileSet(tilemap.Lazy{Count: len(names), Width: isoW, Height: isoH,
		Make: func(id int) *canvas.Canvas { return tileLook(names[id]) }})
	tiles.CellWidth, tiles.CellHeight = tile, tile
	id := func(names ...string) []int {
		ids := make([]int, len(names))
		for i, n := range names {
			ids[i] = tileID[n]
		}
		return ids
	}
	tiles.SetSolid(true, id("fence", "block", "tree0", "anvil")...)
	tiles.SetSolid(true, id(interiorNames[len(interiorFloors):]...)...)
	// Grass, flowers and trees sway each in their own time; water ripples
	// all together, shores and all.
	tiles.Animate(tileID["grass0"], tilemap.Animation{Frames: id("grass0", "grass1", "grass2", "grass3"), FPS: 2, Staggered: true})
	tiles.Animate(tileID["flowers0"], tilemap.Animation{Frames: id("flowers0", "flowers1", "flowers2", "flowers3"), FPS: 2, Staggered: true})
	tiles.Animate(tileID["tree0"], tilemap.Animation{Frames: id("tree0", "tree1", "tree2", "tree3"), FPS: 1.5, Staggered: true})
	// Water is no wall: the farmer swims in it.
	for _, mask := range terrain.Masks() {
		tiles.Animate(shore(mask), tilemap.Animation{Frames: id(shoreName(mask, 0), shoreName(mask, 1), shoreName(mask, 2), shoreName(mask, 3)), FPS: 4})
	}
	// A ripe crop sways, asking to be picked.
	tiles.Animate(tileID["carrot3"], tilemap.Animation{Frames: id("carrot3", "carrot4"), FPS: 2, Staggered: true})
	tiles.Animate(tileID["tomato3"], tilemap.Animation{Frames: id("tomato3", "tomato4"), FPS: 2, Staggered: true})
	return tiles
}

// interiorFloors are the floors of the houses' insides, and interiorNames
// all their tiles, the floors first: what stands on them is solid.
var (
	interiorFloors = []string{"floor", "mat", "rug0", "rug1", "rug2", "rug3"}
	interiorNames  = append(append([]string{}, interiorFloors...),
		"forge", "bed", "bed_blue", "table", "chair", "shelf", "stove", "plant", "chest",
		"wall_n0", "wall_nwin0", "wall_w0", "wall_nw0", "wall_n1", "wall_nwin1", "wall_w1", "wall_nw1",
		"wall_n2", "wall_nwin2", "wall_w2", "wall_nw2", "wall_n3", "wall_nwin3", "wall_w3", "wall_nw3")
)
