package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// clearCell makes the w by h pixels at x, y of img transparent, which
// filling with a transparent colour would not: filling blends.
func clearCell(img *canvas.Canvas, x, y, w, h int) {
	for r := y; r < y+h; r++ {
		row := img.Pixels[r*img.Stride+x : r*img.Stride+x+w]
		for i := range row {
			row[i] = 0
		}
	}
}
