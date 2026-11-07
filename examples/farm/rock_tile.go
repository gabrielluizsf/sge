package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/farm"
)

// rockModels are the models the rocks are made of, by the variant of the
// rock: a small rock and a boulder, seen from two sides. The manifest's
// coast_rocks_01 is a stretch of shore, not a rock one could mine.
var rockModels = [farm.RockVariants]string{"rock_07", "boulder_01", "boulder_01"}

// rockTile is the rock of the tile name: its variant's model, baked, at
// the heading its biome turns it to, standing on its cell, whole. Nil while
// it is not baked, or when name is no rock's.
func (a *realArt) rockTile(name string) *canvas.Canvas {
	v, ok := farm.RockVariant(name)
	if !ok {
		return nil
	}
	turn := v
	for _, c := range farm.RockBiome(name) {
		turn += int(c)
	}
	model := rockModels[v]
	b, ok := a.first(fmt.Sprintf("%s_%d", model, turn%4), fmt.Sprintf("%s_v0_%d", model, turn%2))
	if !ok {
		return nil
	}
	return b.standingAll()
}
