package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

// pictures holds the drawings rendered lately, by their source and size:
// an SVG is parsed and rendered once while it is wanted.
var pictures pictureCache

// picture renders the SVG src, w by h pixels of art, at the game's scale.
func picture(src string, w, h int) *canvas.Canvas {
	key := fmt.Sprint(w, "x", h, src)
	if cv, ok := pictures.get(key); ok {
		return cv
	}
	s, err := sprite.FromSVG(src, w*scale, h*scale)
	if err != nil {
		panic(err)
	}
	pictures.put(key, s.Image)
	return s.Image
}
