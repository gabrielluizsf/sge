package farm

import (
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// penRect is the farm's pen: a fence round the cells it covers, on the
// grass south of the plots, with a gate in the middle of its south side.
var penRect = rect{7, 12, 6, 5}

// penGate is the cell of the fence left open for the player to go in by.
var penGate = town.Cell{Col: penRect.col + penRect.cols/2, Row: penRect.row + penRect.rows - 1}

// penCols and penRows are the size of a pen, fence and all.
const (
	penCols = 6
	penRows = 5
)

// fenceAt is the tile of pen r's fence on cell col, row, -1 off it or at its
// gate.
func fenceAt(r rect, col, row int) int {
	switch {
	case !r.has(col, row) || col == r.col+r.cols/2 && row == r.row+r.rows-1:
		return -1
	case row == r.row || row == r.row+r.rows-1:
		return tileID["fence"]
	case col == r.col || col == r.col+r.cols-1:
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
			switch {
			case col == penGate.Col && row == penGate.Row:
				l.things.Set(col, row, tilemap.Empty)
			case row == r.row || row == r.row+r.rows-1:
				l.things.Set(col, row, tileID["fence"])
			case col == r.col || col == r.col+r.cols-1:
				l.things.Set(col, row, tileID[fenceDown])
			default:
				l.things.Set(col, row, tilemap.Empty)
			}
		}
	}
}
