package sprite

import "github.com/gabrielluizsf/antui/canvas"

// Picture is the sprite's image, made by Make the first time it is asked
// for when it has none.
func (s *Sprite) Picture() *canvas.Canvas {
	s.ready()
	return s.Image
}
