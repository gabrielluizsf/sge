package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// isOpaque reports whether every pixel of img is wholly opaque, so it may
// be copied over what is under it rather than blended.
func isOpaque(img *canvas.Canvas) bool {
	for y := range img.Height {
		for x := range img.Width {
			if img.At(x, y).A() != 255 {
				return false
			}
		}
	}
	return true
}
