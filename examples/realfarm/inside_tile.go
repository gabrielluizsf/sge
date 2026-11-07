package main

import (
	"fmt"
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
	// turn is how many quarter turns on from heading 0 the baked picture is
	// that lies the way the game lays the piece.
	turn int
}

// furnishings are the farm's furniture, by its names for them, each a
// CC0 model of Poly Haven's at its real size: the bed two metres long.
var furnishings = map[string]furnishing{
	"bed_double": {bakes.DoubleBedModel, 2, 2, bakes.DoubleBedTurn}, "bed": {id: "vintage_day_bed", cols: 2, rows: 1}, "bed_blue": {id: "vintage_day_bed", cols: 2, rows: 1},
	"table": {id: "wooden_table_02", cols: 1, rows: 1}, "chair": {id: "painted_wooden_chair_02", cols: 1, rows: 1},
	"shelf": {id: "painted_wooden_shelves", cols: 1, rows: 1}, "plant": {id: "potted_plant_02", cols: 1, rows: 1}, "chest": {id: "treasure_chest", cols: 1, rows: 1},
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
// furniture, drawn from its back cell, its picture reaching forward over
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

// furnitureTile is furnishing f standing on the cells it covers, drawn
// from the back one: its middle half a cell forward, down the screen, for
// each cell more it covers across or deep.
func (a *realArt) furnitureTile(f furnishing) *canvas.Canvas { return a.turned(f, 0) }

// Turned is the piece name turned heading quarter turns, as its model is
// baked at four headings, standing as Tile has it; nil when that heading
// is not baked. A quarter turn swaps the cells it covers across and deep.
func (a *realArt) Turned(name string, heading int) *canvas.Canvas {
	if f, ok := furnishings[name]; ok {
		return a.turned(f, heading)
	}
	return nil
}

// turned is f at heading, or nil if it is not baked at that heading. A
// model of variants side by side is baked at two headings only, the ones
// after those turn it the other way round.
func (a *realArt) turned(f furnishing, heading int) *canvas.Canvas {
	picture := (heading + f.turn) % 4
	b, ok := a.first(fmt.Sprintf("%s_%d", f.id, picture), fmt.Sprintf("%s_v0_%d", f.id, picture))
	if !ok {
		return nil
	}
	cols, rows := f.cols, f.rows
	if heading%2 == 1 {
		cols, rows = rows, cols
	}
	fwd := [2]float64{float64(cols-1) / 2, float64(rows-1) / 2}
	return b.standingOff(int((fwd[0]-fwd[1])*farm.TileW/2), int((fwd[0]+fwd[1])*farm.TileW/4), (cols+rows-2)*farm.TileW/4)
}

// SeatBaked says the sitter is baked with the chair, by bakes.seatOnChair.
func (a *realArt) SeatBaked() bool { return true }
