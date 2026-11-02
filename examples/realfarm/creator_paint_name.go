package main

import (
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// caretBlink is how long the caret is shown, and hidden, in turn.
const caretBlink = 500 * time.Millisecond

// paintName draws the name row: the label, and the name in a field that
// is lit, with a caret, while it has the keyboard.
func (c *creator) paintName(cv *canvas.Canvas, row *row, line rect) {
	c.text(cv, row.label, 'k', [2]int{line.X + 16, line.Y + line.H/2})
	x0 := controlX(line)
	box := rect{x0, line.Y + 8, line.X + line.W - 8 - x0, line.H - 16}
	edge := trackInk
	if c.typing {
		edge = knobInk
	}
	cv.FillRect(canvas.Area{X: box.X - 1, Y: box.Y - 1, Width: box.W + 2, Height: box.H + 2}, edge)
	cv.FillRect(canvas.Area{X: box.X, Y: box.Y, Width: box.W, Height: box.H}, edgeInk)
	x := box.X + 10
	if s := row.text(c.st); s != "" {
		x += c.text(cv, s, 'w', [2]int{x, box.Y + box.H/2})
	}
	if c.typing && time.Now().UnixMilli()/caretBlink.Milliseconds()%2 == 0 {
		cv.FillRect(canvas.Area{X: x + 2, Y: box.Y + 6, Width: 2, Height: box.H - 12}, knobInk)
	}
}
