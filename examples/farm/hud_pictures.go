package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/label"
)

// hudPictures holds the HUD's text and panels by what they show, so each
// is drawn once, not every frame.
var hudPictures = map[string]*canvas.Canvas{}

// outline is how many pixels round text over the world is outlined.
const outline = 2

// text is s in the HUD's face, in the palette's ink: outlined round, but
// for the panels' cocoa.
func text(s string, ink rune) *canvas.Canvas {
	key := fmt.Sprint("t", string(ink), s)
	if cv, ok := hudPictures[key]; ok {
		return cv
	}
	st := label.Style{Ink: inks[ink]}
	if outlined(ink) {
		st.Outline, st.Width = inkOutline, outline
	}
	cv := label.New(hudFace, s, st)
	hudPictures[key] = cv
	return cv
}

// panel is a HUD panel of w by h pixels of art, gold when selected.
func panel(w, h int, selected bool) *canvas.Canvas {
	key := fmt.Sprint("p", w, "x", h, selected)
	if cv, ok := hudPictures[key]; ok {
		return cv
	}
	cv := picture(art.Panel(w, h, selected), w, h)
	hudPictures[key] = cv
	return cv
}
