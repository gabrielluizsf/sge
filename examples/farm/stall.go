package main

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// The farm's stall's cell; a shop is open to one standing in front of a
// stall.
const stallCol, stallRow = 11, 3

// farmStall is the farm's stall, and home the player's house: its two by
// two cells from homeCorner, and its doorstep.
var (
	farmStall  = town.Cell{Col: stallCol, Row: stallRow}
	homeCorner = town.Cell{Col: 16, Row: 1}
	homeDoor   = town.Cell{Col: 16, Row: 3}
)

// stallBy is the stall o stands in front of, the farm's or a town's.
func (f *farm) stallBy(o *game.Object) (town.Cell, bool) {
	col, row := f.cellOf(feet(o))
	stalls := []town.Cell{farmStall}
	for _, t := range f.towns.Near(col, row, town.Reach+4) {
		stalls = append(stalls, t.Stall)
	}
	for _, s := range stalls {
		if row == s.Row+1 && col >= s.Col-1 && col <= s.Col+1 {
			return s, true
		}
	}
	return town.Cell{}, false
}

// atStall reports whether o stands in front of a stall's counter.
func (f *farm) atStall(o *game.Object) bool {
	_, ok := f.stallBy(o)
	return ok
}
