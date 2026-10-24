package main

import (
	"strings"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// layers are the farm's maps, from the ground up: the ground, its pond,
// what grows, what stands on it, and what the farmer walks under. The
// ground is grass under everything; the rest are transparent around what
// they show, so the grass beneath sways.
type layers struct {
	ground, water, crops, things, over *tilemap.Map
}

// The layers' places in the drawing order, and the farmer's and the
// HUD's among them. The endless world around the farm is drawn at the
// same places.
const (
	zGround = iota
	zWater
	zCrops
	zThings
	zFarmer
	zOver
	zHUD
)

// newLayers lays worldMap out in tiles and adds the layers to g.
func newLayers(g *game.Game, tiles *tilemap.TileSet) layers {
	grass := tileID["grass0"]
	ground := map[rune]int{'G': tileID["flowers0"], 'p': tileID["path"], 'd': tileID["soil_dry"]}
	for _, c := range "gwfTtRrCc" {
		ground[c] = grass
	}
	things := map[rune]int{'f': tileID["fence"], 'T': tileID["tree_trunk"], 'C': tileID["stall_counter_l"], 'c': tileID["stall_counter_r"]}
	over := map[rune]int{'t': tileID["tree_crown0"], 'R': tileID["stall_roof_l"], 'r': tileID["stall_roof_r"]}
	l := layers{
		ground: parse("ground", tiles, ground),
		water:  parse("water", tiles, nil),
		crops:  parse("crops", tiles, nil),
		things: parse("things", tiles, things),
		over:   parse("over", tiles, over),
	}
	// The pond's shores join up as its cells lie.
	rows := strings.Split(strings.Trim(worldMap, "\n"), "\n")
	wet := func(c, r int) bool { return r >= 0 && r < len(rows) && c >= 0 && c < len(rows[r]) && rows[r][c] == 'w' }
	for r, line := range rows {
		for c := range line {
			if wet(c, r) {
				l.water.Set(c, r, shore(terrain.Mask(wet, c, r)))
			}
		}
	}
	for i, m := range []*tilemap.Map{l.ground, l.water, l.crops, l.things, l.over} {
		m.Z = []int{zGround, zWater, zCrops, zThings, zOver}[i]
		g.Add(m.Object)
	}
	return l
}

func parse(name string, tiles *tilemap.TileSet, legend map[rune]int) *tilemap.Map {
	m, err := tilemap.Parse(name, tiles, worldMap, legend)
	if err != nil {
		panic(err)
	}
	return m
}
