package farm

import "github.com/gabrielluizsf/antui/canvas"

// pieceDrop is how far below the back cell's bottom corner a piece of cols
// by rows cells reaches: down to its front cell's. The whole picture of a
// piece ends there, its middle over the back cell's corner.
func pieceDrop(cols, rows int) int { return (cols + rows - 2) * isoW / 4 }

// sliceOf is the part of full, the picture of a piece cols by rows cells
// large drawn from its back cell, that belongs to cell at (its column and
// row in the piece), ready to be drawn at that cell: placed by its bottom
// middle on the cell's bottom corner, as every tile is.
//
// A pixel belongs to a cell by where the ground under it falls on the
// piece's footprint, seen from the corner: the cells' edges are the
// diagonals x+2y and 2y-x of the picture, taken from the back cell's bottom
// corner. What stands tall is thus cut by the same diagonals, and whatever
// reaches beyond the piece's cells belongs to the nearest cell, so every
// pixel belongs to one cell and the parts, drawn one over another in the
// order of their cells' depths, make the picture again.
func sliceOf(full *canvas.Canvas, cols, rows int, at [2]int) *canvas.Canvas {
	sx, sy := (at[0]-at[1])*isoW/2, (at[0]+at[1])*isoW/4
	w, h := full.Width+2*abs(sx), full.Height+sy
	out, err := canvas.NewCanvas(w, h)
	if err != nil {
		return full
	}
	for oy := range h {
		for ox := range w {
			// x, y: the pixel from the bottom corner of the back cell.
			x, y := ox-w/2+sx, oy-h+sy
			if cellOf(x, y, cols, rows) != at {
				continue
			}
			if c := full.At(x+full.Width/2, y+full.Height-pieceDrop(cols, rows)); c.A() != 0 {
				out.Pixels[oy*out.Stride+ox] = c
			}
		}
	}
	return out
}

// cellOf is the cell of a piece cols by rows large that the pixel x, y,
// from the bottom corner of the back cell, belongs to. A pixel lower than
// its cell's bottom corner, which that cell's frame could not draw, goes to
// the next cell forward, on the side of the piece it lies.
func cellOf(x, y, cols, rows int) [2]int {
	i := max(0, min(cols-1, floorDiv(x+2*y, isoW)+1))
	j := max(0, min(rows-1, floorDiv(2*y-x, isoW)+1))
	for y >= (i+j)*isoW/4 {
		switch {
		case i < cols-1 && (x >= (i-j)*isoW/2 || j == rows-1):
			i++
		case j < rows-1:
			j++
		default:
			return [2]int{i, j}
		}
	}
	return [2]int{i, j}
}
