package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/human"
)

// dropperW is the width of the eyedropper's button in a code row.
const dropperW = 72

// paintCode draws a colour's code row: the label, the code in a field
// that is lit while it has the keyboard, and the eyedropper's button.
func (c *creator) paintCode(cv *canvas.Canvas, row *row, line rect) {
	c.text(cv, i18n.T(row.label), 'k', [2]int{line.X + 16, line.Y + line.H/2})
	x0 := controlX(line)
	box := rect{x0, line.Y + 8, line.X + line.W - 8 - x0 - dropperW - 6, line.H - 16}
	text := codeOf(*row.field(c.st))
	if text == "" {
		text = i18n.T("Original")
	}
	lit := c.typing && c.at() == row
	if lit {
		text = c.code
	}
	c.paintField(cv, fieldView{text, box, lit})
	btn := rect{box.X + box.W + 6, box.Y, dropperW, box.H}
	cv.FillRect(canvas.Area{X: btn.X - 1, Y: btn.Y - 1, Width: btn.W + 2, Height: btn.H + 2}, trackInk)
	cv.FillRect(canvas.Area{X: btn.X, Y: btn.Y, Width: btn.W, Height: btn.H}, edgeInk)
	c.centred(cv, i18n.T("Pick"), 'y', btn)
	idx := len(c.hits) - 1
	c.hits = append(c.hits, hit{at: btn, kind: hitDropper, cat: c.hits[idx].cat, row: c.hits[idx].row})
}

// paintNote draws the message of the code field, or the eyedropper's
// colour under the mouse, along the foot of the rows' panel.
func (c *creator) paintNote(cv *canvas.Canvas, r rect) {
	foot := rect{r.X + 12, r.Y + r.H - 56, r.W - 24, 44}
	if d := c.drop; d != nil {
		if d.loaded {
			col := human.RGB(float64(d.at.R)/255, float64(d.at.G)/255, float64(d.at.B)/255)
			fillColour(cv, rect{foot.X + 8, foot.Y + 8, 28, 28}, col)
			c.text(cv, codeOf(col), 'k', [2]int{foot.X + 48, foot.Y + 22})
		}
		c.centred(cv, i18n.T("Click any pixel of the screen. Esc: cancel"), 'k', rect{foot.X + 150, foot.Y, foot.W - 150, foot.H})
		return
	}
	if c.note != "" {
		c.centred(cv, c.note, 'k', foot)
	}
}
