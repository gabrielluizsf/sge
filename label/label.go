// Package label draws text as a picture for a game's HUD or its world: a
// face's glyphs in an ink, smooth at their edges, outlined round so they
// read on any ground, the picture transparent everywhere else.
package label

import "github.com/gabrielluizsf/antui/canvas"

// Style is how a label is drawn: its Ink, and an Outline Width pixels
// round it; a Width of 0 draws no outline.
type Style struct {
	Ink, Outline canvas.Color
	Width        int
}
