package main

import "github.com/gabrielluizsf/sge/game"

// The stall's counter's left cell; the shop is open to a farmer standing
// in front of it.
const stallCol, stallRow = 11, 3

// atStall reports whether o stands in front of the stall's counter.
func (f *farm) atStall(o *game.Object) bool {
	col, row, _ := f.ground.CellAt(feet(o).X, feet(o).Y)
	return row == stallRow+1 && col >= stallCol-1 && col <= stallCol+2
}
