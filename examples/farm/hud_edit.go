package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// editView draws the arranging of furniture over the world: the piece
// about to be put down as a ghost on its cells, green where it fits and
// red where it does not, or the cells of the piece under the cursor, gold;
// and the keys under the navbar's toasts.
func (h *hud) editView(cv *canvas.Canvas) {
	f := h.f
	e := f.edit
	if e.pen != nil {
		r, ok := f.penGhost()
		tint := ghostRed
		if ok {
			tint = ghostGreen
		}
		h.outlineRect(cv, r, tint)
	} else if p := f.penUnder(e.col, e.row); p != nil && e.hold == nil {
		h.outlineRect(cv, p.r, ghostGold)
	} else if p, ok, shown := f.ghost(); shown {
		tint := ghostRed
		if ok {
			tint = ghostGreen
		}
		h.ghostPiece(cv, p, tint)
	} else if p, ok := f.pieceUnder(e.col, e.row); ok {
		h.outline(cv, p, ghostGold)
	} else {
		h.outline(cv, newPiece("chair", e.col, e.row, 0), ghostGold)
	}
	h.editHints(cv)
}

// The ghost's tints.
var (
	ghostGreen = canvas.RGB(120, 255, 140)
	ghostRed   = canvas.RGB(255, 90, 80)
	ghostGold  = canvas.RGB(255, 214, 90)
)

// ghostPiece draws piece p translucent, tinted, standing on its cells, and
// the diamonds of its cells outlined.
func (h *hud) ghostPiece(cv *canvas.Canvas, p piece, tint canvas.Color) {
	f := h.f
	img := f.tinted(fmt.Sprint("ghost ", p.name, p.heading, tint), func() *canvas.Canvas { return pieceFull(f.art, p.name, p.heading) }, tint, 0.6)
	at := h.farmer.Game().ToScreen(f.ground.CellCenter(p.col, p.row))
	cols, rows := footprint(f.art, p.name, p.heading)
	cv.Blit(int(at.X)-img.Width/2, int(at.Y)+tile/2+pieceDrop(cols, rows)-img.Height, img)
	h.outline(cv, p, tint)
}

// outline draws the diamond of each cell piece p covers, tinted.
func (h *hud) outline(cv *canvas.Canvas, p piece, tint canvas.Color) {
	f := h.f
	cur := f.tinted(fmt.Sprint("cursor ", tint), f.art.Cursor, tint, 1)
	s := p.cells(f.art)
	for r := s.r0; r < s.r0+s.rows; r++ {
		for c := s.c0; c < s.c0+s.cols; c++ {
			at := h.farmer.Game().ToScreen(f.ground.CellCenter(c, r))
			cv.Blit(int(at.X)-isoW/2, int(at.Y)-tile/2, cur)
		}
	}
}

// outlineRect draws the diamond of each cell of r, tinted.
func (h *hud) outlineRect(cv *canvas.Canvas, r rect, tint canvas.Color) {
	f := h.f
	cur := f.tinted(fmt.Sprint("cursor ", tint), f.art.Cursor, tint, 1)
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			at := h.farmer.Game().ToScreen(f.ground.CellCenter(col, row))
			cv.Blit(int(at.X)-isoW/2, int(at.Y)-tile/2, cur)
		}
	}
}

// editHints says what is being placed and the keys, under the toasts.
func (h *hud) editHints(cv *canvas.Canvas) {
	f := h.f
	e := f.edit
	status := i18n.T("Arranging: nothing chosen")
	switch {
	case e.pen != nil:
		status = i18n.T("Holding: %s", i18n.T("Pen"))
	case f.penUnder(e.col, e.row) != nil:
		status = i18n.T("Enter: move this pen")
	case e.hold != nil:
		status = i18n.T("Holding: %s", i18n.T(pieceLabels[e.hold.name]))
	case e.item != "":
		status = i18n.T("Place: %s (Q/E: change)", i18n.T(pieceLabels[e.item]))
	}
	keys := i18n.T("Arrows: move   Enter: lift/place   R: turn   X: store   Esc/B: done")
	y := (navBottom + 4) * scale
	for _, line := range []struct {
		s   string
		ink rune
	}{{status, 'y'}, {keys, 'w'}} {
		t := h.text(line.s, line.ink)
		cv.Blit((cv.Width-t.Width)/2, y, t)
		y += t.Height + 2*scale
	}
}
