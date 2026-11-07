package farm

import (
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// penRect is the farm's pen: a fence round the cells it covers, on the
// grass south of the plots, with a gate in the middle of its south side.
var penRect = rect{5, 12, penCols, penRows}

// penGate is the cell of the fence left open for the player to go in by.
var penGate = town.Cell{Col: penRect.col + penRect.cols/2, Row: penRect.row + penRect.rows - 1}

// penCols and penRows are the size of a pen, fence and all.
const (
	penCols = 8
	penRows = 6
)

// fenceAt is the tile of pen r's fence on cell col, row, -1 off it or at its
// gate: a corner at each corner, so its sides meet, and runs between.
func fenceAt(r rect, col, row int) int {
	north, south := row == r.row, row == r.row+r.rows-1
	west, east := col == r.col, col == r.col+r.cols-1
	switch {
	case !r.has(col, row) || south && col == r.col+r.cols/2:
		return -1
	case north && west:
		return tileID[fenceCorners[0]]
	case north && east:
		return tileID[fenceCorners[1]]
	case south && west:
		return tileID[fenceCorners[2]]
	case south && east:
		return tileID[fenceCorners[3]]
	case north || south:
		return tileID["fence"]
	case west || east:
		return tileID[fenceDown]
	}
	return -1
}

// inside is the cells of the pen's field, within its fence.
func (r rect) inside() rect { return rect{r.col + 1, r.row + 1, r.cols - 2, r.rows - 2} }

// applyPen fences the pen on the things of l, with no tree in it.
func applyPen(l layers) {
	r := penRect
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			if id := fenceAt(r, col, row); id >= 0 {
				l.things.Set(col, row, id)
				continue
			}
			l.things.Set(col, row, tilemap.Empty)
		}
	}
}
