package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vignette"
)

// hud draws what the player needs to see over the world, fixed to the
// screen, as farm games lay it out: the money, the time and who the
// player is and how old at the top left, what they can hold in a bar at
// the bottom, the cell they work marked, bubbles over those who want something, words rising from what
// happens, news in the middle, and a question, the shop or a
// conversation over it all when open.
type hud struct {
	f      *farm
	farmer *game.Object
	p      *farmer
	// round is the player's portrait as a circle, made from roundSrc at
	// roundD pixels, kept until either changes.
	round    *canvas.Canvas
	roundSrc *canvas.Canvas
	roundD   int
	// t is the seconds the HUD has been shown, for what pulses.
	t float64
	// edge is the red of a blow at full strength, for the size of the screen.
	edge *vignette.Vignette
}

func (h *hud) Update(_ *game.Object, dt float64) { h.t += dt }

func (h *hud) Paint(o *game.Object, cv *canvas.Canvas) {
	menu := h.f.shop.open || h.f.talk != nil || h.f.question != nil || h.f.ui.open != playing || h.f.pack.Open
	if !menu {
		h.cursor(cv)
	}
	h.bubbles(cv)
	if !menu {
		h.cropLabel(cv)
	}
	for _, w := range h.f.floats {
		t := h.text(w.text, w.ink)
		at := o.Game().ToScreen(w.at)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-t.Height, t)
	}
	if h.f.edit != nil && !menu {
		h.editView(cv)
	}
	h.hurtEdge(cv)
	h.navbar(cv)
	h.hotbar(cv)
	h.notices(cv)
	if !menu {
		h.hints(cv)
	}
	switch {
	case h.f.question != nil:
		h.questionPanel(cv)
	case h.f.shop.open:
		h.shop(cv)
	case h.f.talk != nil:
		h.talk(cv)
	}
	h.f.ui.paint(cv)
}

// cursor marks the cell the farmer works, a plot, or a tree with the axe
// in hand: its diamond, outlined.
func (h *hud) cursor(cv *canvas.Canvas) {
	col, row := h.p.useCell(h.farmer)
	f := h.f
	if h.p.swims || f.room != nil || !f.marks(col, row) {
		return
	}
	at := h.farmer.Game().ToScreen(f.ground.CellCenter(col, row))
	cv.Blit(int(at.X)-isoW/2, int(at.Y)-tile/2, h.f.art.Cursor())
}

// rowH is the height, in pixels of art, of a line of the HUD's panels.
const rowH = 16
