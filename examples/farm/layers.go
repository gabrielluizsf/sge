package main

import (
	"strings"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
)

// layers are the farm's maps, from the ground up: the ground, its pond,
// what grows and what stands on it. The ground is grass under everything;
// the rest are transparent around what they show, so the grass beneath
// sways. What grows and what stands are Sorted: seen from the corner, each
// stands among the farmer and the stall, nearer over farther.
type layers struct {
	ground, water, crops, things *tilemap.Map
}

// The layers' places in the drawing order, and the HUD's over them. The
// farmer and the stall stand among what stands on the ground. The endless
// world around the farm is drawn at the same places.
const (
	zGround = iota
	zWater
	zThings
	zHUD
)

// newLayers lays worldMap out in tiles and adds the layers to g.
func newLayers(g *game.Game, tiles *tilemap.TileSet) layers {
	grass := tileID["grass0"]
	ground := map[rune]int{'G': tileID["flowers0"], 'p': tileID["path"], 'd': tileID["soil_dry"]}
	for _, c := range "gwfTtC" {
		ground[c] = grass
	}
	// The stall is an object on its cell (see newStall); it is solid.
	things := map[rune]int{'f': tileID["fence"], 'T': tileID["tree0"], 'C': tileID["block"]}
	l := layers{
		ground: parse("ground", tiles, ground),
		water:  parse("water", tiles, nil),
		crops:  parse("crops", tiles, nil),
		things: parse("things", tiles, things),
	}
	l.crops.Sorted, l.things.Sorted = true, true
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
	for i, m := range []*tilemap.Map{l.ground, l.water, l.crops, l.things} {
		m.Z = []int{zGround, zWater, zThings, zThings}[i]
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
