package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// copyTile puts tile into img at x, y in place of what was there, its
// transparency too: a map's picture is its own, so its pixels are written
// directly, with no blending.
func copyTile(img *canvas.Canvas, x, y int, tile *canvas.Canvas) {
	for r := range tile.Height {
		copy(img.Pixels[(y+r)*img.Stride+x:(y+r)*img.Stride+x+tile.Width], tile.Pixels[r*tile.Stride:r*tile.Stride+tile.Width])
	}
}
