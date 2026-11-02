package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/human"
)

// The creation screen's measures, in screen pixels.
const (
	edgeGap  = 32  // from the window's edge to a panel
	listW    = 300 // the categories' panel
	rowsW    = 500 // the rows' panel
	headH    = 64  // the panels' title
	catH     = 48  // a category's line
	rowH     = 64  // a row
	arrowW   = 40
	trackH   = 8
	swatchSz = 22
)

// Colours of the controls: parchment track and gold fill, as the panels'.
var (
	trackInk = canvas.RGB(0x4a, 0x3c, 0x25)
	fillInk  = canvas.RGB(0xb8, 0x92, 0x4a)
	knobInk  = canvas.RGB(0xf0, 0xc5, 0x5e)
	edgeInk  = canvas.RGB(0x14, 0x10, 0x0f)
)

// dimmed is the shade laid over the world behind the screen.
var dimmed = canvas.RGBA(canvas.RGBAParams{R: 10, G: 12, B: 24, A: 150})

// Paint draws the screen over the world, which goes on behind it, and
// places what the mouse acts on.
func (c *creator) paint(cv *canvas.Canvas) {
	cv.FillRect(canvas.Area{Width: cv.Width, Height: cv.Height}, dimmed)
	c.hits = c.hits[:0]
	h := cv.Height - 2*edgeGap
	c.paintList(cv, rect{edgeGap, edgeGap, listW, h})
	c.paintRows(cv, rect{cv.Width - edgeGap - rowsW, edgeGap, rowsW, h})
	c.paintStage(cv, rect{2*edgeGap + listW, edgeGap, cv.Width - 4*edgeGap - listW - rowsW, h})
}

// panel draws a panel over r, lit if chosen.
func (c *creator) panel(cv *canvas.Canvas, r rect, chosen bool) {
	cv.Blit(r.X, r.Y, c.art.Panel(r.W/farm.Scale, r.H/farm.Scale, chosen))
}

// text draws s in ink at x, y, centred on y; the width it took.
func (c *creator) text(cv *canvas.Canvas, s string, ink rune, at [2]int) int {
	t := c.art.Text(s, ink)
	cv.Blit(at[0], at[1]-t.Height/2, t)
	return t.Width
}

// centred draws s in ink centred in r.
func (c *creator) centred(cv *canvas.Canvas, s string, ink rune, r rect) {
	t := c.art.Text(s, ink)
	cv.Blit(r.X+(r.W-t.Width)/2, r.Y+(r.H-t.Height)/2, t)
}

// paintList draws the categories.
func (c *creator) paintList(cv *canvas.Canvas, r rect) {
	c.panel(cv, r, false)
	c.centred(cv, "Criar personagem", 'y', rect{r.X, r.Y, r.W, headH})
	for i, cat := range c.cats {
		line := rect{r.X + 12, r.Y + headH + i*catH, r.W - 24, catH - 4}
		chosen := i == c.cat
		if chosen {
			c.panel(cv, line, true)
		}
		c.hits = append(c.hits, hit{at: line, kind: hitCategory, cat: i})
		ink := 'k'
		if chosen {
			ink = 'y'
		}
		c.centred(cv, cat.name, ink, line)
	}
}

// paintRows draws the rows of the category chosen.
func (c *creator) paintRows(cv *canvas.Canvas, r rect) {
	c.panel(cv, r, false)
	cat := c.cats[c.cat]
	c.centred(cv, cat.name, 'y', rect{r.X, r.Y, r.W, headH})
	y := r.Y + headH
	for _, i := range c.shown(c.cat) {
		line := rect{r.X + 12, y, r.W - 24, rowH - 4}
		row := &cat.rows[i]
		chosen := i == c.row
		if chosen {
			c.panel(cv, line, true)
		}
		c.hits = append(c.hits, hit{at: line, kind: hitRow, cat: c.cat, row: i})
		switch row.kind {
		case buttonRow:
			c.centred(cv, row.label, 'y', line)
		case choiceRow:
			c.paintChoice(cv, row, line)
		case sliderRow:
			c.paintSlider(cv, row, line)
		case textRow:
			c.paintName(cv, row, line)
		}
		y += rowH
	}
}

// controlX is where a row's control begins, the label before it.
func controlX(line rect) int { return line.X + line.W*4/10 }

// paintChoice draws a choice row: the label, and « value » .
func (c *creator) paintChoice(cv *canvas.Canvas, row *row, line rect) {
	c.text(cv, row.label, 'k', [2]int{line.X + 16, line.Y + line.H/2})
	x0 := controlX(line)
	w := line.X + line.W - 8 - x0
	less := rect{x0, line.Y, arrowW, line.H}
	more := rect{x0 + w - arrowW, line.Y, arrowW, line.H}
	c.centred(cv, "«", 'y', less)
	c.centred(cv, "»", 'y', more)
	idx := len(c.hits) - 1 // the row's own hit, which the arrows' go over
	c.hits = append(c.hits,
		hit{at: less, kind: hitLess, cat: c.hits[idx].cat, row: c.hits[idx].row},
		hit{at: more, kind: hitMore, cat: c.hits[idx].cat, row: c.hits[idx].row})
	mid := rect{x0 + arrowW, line.Y, w - 2*arrowW, line.H}
	t := c.art.Text(row.text(c.st), 'k')
	sw := 0
	if row.swatch != nil {
		if col := row.swatch(c.st); col.Set {
			sw = swatchSz + 8
		}
	}
	x := mid.X + (mid.W-t.Width-sw)/2
	if sw > 0 {
		fillColour(cv, rect{x, line.Y + (line.H-swatchSz)/2, swatchSz, swatchSz}, row.swatch(c.st))
	}
	cv.Blit(x+sw, line.Y+(line.H-t.Height)/2, t)
}

// fillColour draws the colour col in r, edged dark.
func fillColour(cv *canvas.Canvas, r rect, col human.Colour) {
	cv.FillRect(canvas.Area{X: r.X - 1, Y: r.Y - 1, Width: r.W + 2, Height: r.H + 2}, edgeInk)
	cv.FillRect(canvas.Area{X: r.X, Y: r.Y, Width: r.W, Height: r.H},
		canvas.RGB(uint8(col.R*255+0.5), uint8(col.G*255+0.5), uint8(col.B*255+0.5)))
}

// paintSlider draws a slider row: the label above its track and knob.
func (c *creator) paintSlider(cv *canvas.Canvas, row *row, line rect) {
	c.text(cv, row.label, 'k', [2]int{line.X + 16, line.Y + line.H/2})
	x0 := controlX(line)
	track := rect{x0 + 8, line.Y + (line.H-trackH)/2, line.X + line.W - 16 - x0 - 8, trackH}
	f := min(max(row.slider.frac(c.st), 0), 1)
	cv.FillRect(canvas.Area{X: track.X, Y: track.Y, Width: track.W, Height: track.H}, trackInk)
	cv.FillRect(canvas.Area{X: track.X, Y: track.Y, Width: int(f * float64(track.W)), Height: track.H}, fillInk)
	k := rect{track.X + int(f*float64(track.W)) - 7, track.Y - 6, 14, trackH + 12}
	cv.FillRect(canvas.Area{X: k.X - 1, Y: k.Y - 1, Width: k.W + 2, Height: k.H + 2}, edgeInk)
	cv.FillRect(canvas.Area{X: k.X, Y: k.Y, Width: k.W, Height: k.H}, knobInk)
	idx := len(c.hits) - 1
	c.hits = append(c.hits, hit{at: rect{track.X, line.Y, track.W, line.H}, kind: hitTrack, cat: c.hits[idx].cat, row: c.hits[idx].row})
}

// paintStage draws the person, live, in r.
func (c *creator) paintStage(cv *canvas.Canvas, r rect) {
	c.stageAt = r
	c.panel(cv, r, false)
	inner := c.inner()
	if img := c.stageImage(inner); img != nil {
		cv.Blit(inner.X, inner.Y, img)
	} else {
		c.centred(cv, "Preparando o personagem...", 'k', r)
	}
	if c.live != nil && c.live.Busy() && c.stage.Ready() {
		c.centred(cv, "Atualizando...", 'k', rect{r.X, r.Y + r.H - 56, r.W, 40})
	}
	c.centred(cv, "Arraste para girar", 'k', rect{r.X, r.Y + stageInset + 4, r.W, 32})
}

// stageInset is the margin between the stage's panel and the picture.
const stageInset = 14
