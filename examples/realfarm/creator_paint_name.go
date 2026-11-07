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
	c.paintField(cv, fieldView{row.text(c.st), rect{controlX(line), line.Y + 8, line.X + line.W - 8 - controlX(line), line.H - 16}, c.typing})
}

// fieldView is what a text field shows: its text, where, and whether it
// has the keyboard.
type fieldView struct {
	text string
	box  rect
	lit  bool
}

// paintField draws text in a field over box that is lit, with a caret,
// while it has the keyboard.
func (c *creator) paintField(cv *canvas.Canvas, v fieldView) {
	box, s := v.box, v.text
	edge := trackInk
	if v.lit {
		edge = knobInk
	}
	cv.FillRect(canvas.Area{X: box.X - 1, Y: box.Y - 1, Width: box.W + 2, Height: box.H + 2}, edge)
	cv.FillRect(canvas.Area{X: box.X, Y: box.Y, Width: box.W, Height: box.H}, edgeInk)
	x := box.X + 10
	if s != "" {
		x += c.text(cv, s, 'w', [2]int{x, box.Y + box.H/2})
	}
	if v.lit && time.Now().UnixMilli()/caretBlink.Milliseconds()%2 == 0 {
		cv.FillRect(canvas.Area{X: x + 2, Y: box.Y + 6, Width: 2, Height: box.H - 12}, knobInk)
	}
}
