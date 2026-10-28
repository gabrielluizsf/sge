package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
)

// hudPictures holds the HUD's text and panels by what they show, so the
// SVG of each is written and rendered once, not every frame.
var hudPictures = map[string]*canvas.Canvas{}

// text is s in the HUD's letters, in the palette colour ink.
func text(s string, ink rune) *canvas.Canvas {
	key := fmt.Sprint("t", string(ink), s)
	if cv, ok := hudPictures[key]; ok {
		return cv
	}
	cv := picture(art.Text(s, ink), art.TextWidth(s), 8)
	hudPictures[key] = cv
	return cv
}

// panel is a HUD panel of w by h pixels of art, yellow when selected.
func panel(w, h int, selected bool) *canvas.Canvas {
	key := fmt.Sprint("p", w, "x", h, selected)
	if cv, ok := hudPictures[key]; ok {
		return cv
	}
	cv := picture(art.Panel(w, h, selected), w, h)
	hudPictures[key] = cv
	return cv
}
