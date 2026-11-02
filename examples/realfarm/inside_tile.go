package main

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// furnishing is a model a house's inside is furnished with: its asset,
// and how many cells across and deep it covers.
type furnishing struct {
	id         string
	cols, rows int
}

// furnishings are the farm's furniture, by its names for them, each a
// CC0 model of Poly Haven's at its real size: the bed two metres long.
var furnishings = map[string]furnishing{
	"bed": {"vintage_day_bed", 2, 1}, "bed_blue": {"vintage_day_bed", 2, 1},
	"table": {"wooden_table_02", 1, 1}, "chair": {"painted_wooden_chair_02", 1, 1},
	"shelf": {"painted_wooden_shelves", 1, 1}, "plant": {"potted_plant_02", 1, 1}, "chest": {"treasure_chest", 1, 1},
}

// Furniture is how many cells across and deep the piece name covers.
func (a *realArt) Furniture(name string) (cols, rows int) {
	if f, ok := furnishings[name]; ok {
		return f.cols, f.rows
	}
	return 1, 1
}

// insideTile is the piece of a house's inside name: a floor, a rug, the
// doormat or a back wall, baked from package building, or a piece of
// furniture, drawn from its front cell, its picture reaching back over
// the others it covers. Nil when name is none of them or is not baked.
func (a *realArt) insideTile(name string) *canvas.Canvas {
	if f, ok := furnishings[name]; ok {
		return a.furnitureTile(f)
	}
	if !isInside(name) {
		return nil
	}
	b, ok := a.bake(bakes.InsideName(name))
	if !ok {
		return nil
	}
	return b.standing() // a wall stands a storey tall
}

// isInside reports whether name is one of building.Insides.
func isInside(name string) bool {
	return name == "floor" || name == "mat" || strings.HasPrefix(name, "rug") || strings.HasPrefix(name, "wall_")
}

// furnitureTile is furnishing f standing on the cells it covers, from the
// front one: its middle half a cell back, up the screen, for each cell
// more it covers across or deep.
func (a *realArt) furnitureTile(f furnishing) *canvas.Canvas {
	b, ok := a.first(f.id+"_0", f.id+"_v0_0")
	if !ok {
		return nil
	}
	back := [2]float64{float64(f.cols-1) / 2, float64(f.rows-1) / 2}
	return b.standingOff(-int((back[0]-back[1])*farm.TileW/2), -int((back[0]+back[1])*farm.TileW/4))
}
