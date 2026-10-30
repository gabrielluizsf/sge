package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/label"
)

// outline is how many pixels round text over the world is outlined.
const outline = 2

// Text is s in the HUD's face, in the palette's ink: outlined round, but
// for the panels' cocoa.
func (a *svgArt) Text(s string, ink rune) *canvas.Canvas {
	key := fmt.Sprint("t", string(ink), s)
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	st := label.Style{Ink: inks[ink]}
	if outlined(ink) {
		st.Outline, st.Width = inkOutline, outline
	}
	cv := label.New(hudFace, s, st)
	a.hud[key] = cv
	return cv
}

// Panel is a HUD panel of w by h pixels of art, gold when selected.
func (a *svgArt) Panel(w, h int, selected bool) *canvas.Canvas {
	key := fmt.Sprint("p", w, "x", h, selected)
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	cv := picture(art.Panel(w, h, selected), w, h)
	a.hud[key] = cv
	return cv
}
