package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

// pictures holds every drawing rendered so far, by its source and size:
// an SVG is parsed and rendered once.
var pictures = map[string]*canvas.Canvas{}

// picture renders the SVG src, w by h pixels of art, at the game's scale.
func picture(src string, w, h int) *canvas.Canvas {
	key := fmt.Sprint(w, "x", h, src)
	if cv, ok := pictures[key]; ok {
		return cv
	}
	s, err := sprite.FromSVG(src, w*scale, h*scale)
	if err != nil {
		panic(err)
	}
	pictures[key] = s.Image
	return s.Image
}
