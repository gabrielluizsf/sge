package backpack

import "github.com/gabrielluizsf/antui/canvas"

// Art is how the game draws the pack: its panels, which are w by h pixels
// of its art and brighter when selected; the picture of an item's Icon,
// 16 by 16 pixels of art; and text in an ink, 'k' on a panel.
type Art interface {
	Panel(w, h int, selected bool) *canvas.Canvas
	Icon(name string) *canvas.Canvas
	Text(s string, ink rune) *canvas.Canvas
}
