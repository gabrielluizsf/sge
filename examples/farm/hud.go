package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
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
}

func (h *hud) Update(*game.Object, float64) {}

func (h *hud) Paint(o *game.Object, cv *canvas.Canvas) {
	menu := h.f.shop.open || h.f.talk != nil || h.f.question != nil
	if !menu {
		h.cursor(cv)
	}
	h.bubbles(cv)
	for _, w := range h.f.floats {
		t := h.text(w.text, w.ink)
		at := o.Game().ToScreen(w.at)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-t.Height, t)
	}
	h.money(cv)
	h.time(cv)
	h.age(cv)
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
}

// cursor marks the cell the farmer works, a plot, or a tree with the axe
// in hand: its diamond, outlined.
func (h *hud) cursor(cv *canvas.Canvas) {
	col, row := h.p.target(h.farmer)
	f := h.f
	tree := f.slot == slotAxe && f.thingAt(col, row) == tileID["tree0"]
	if h.p.swims || f.room != nil || f.plotAt(col, row) == nil && !tree {
		return
	}
	at := h.farmer.Game().ToScreen(f.ground.CellCenter(col, row))
	cv.Blit(int(at.X)-isoW/2, int(at.Y)-tile/2, h.f.art.Cursor())
}

// money is the purse at the top left: a coin and what is in it.
func (h *hud) money(cv *canvas.Canvas) {
	t := h.text(fmt.Sprint("$", h.f.money), 'k')
	cv.Blit(8, 8, h.panel(max(52, t.Width/scale+24), rowH, false))
	cv.Blit(8-scale, 8, h.icon("icon_coin"))
	cv.Blit(8+18*scale, 8+(rowH*scale-t.Height)/2, t)
}

// rowH is the height, in pixels of art, of a line of the HUD's panels.
const rowH = 16

// time is the day and the hour under the purse, red once it is late.
func (h *hud) time(cv *canvas.Canvas) {
	c := h.f.clock
	ink := 'k'
	if c.Between(22, wakeAt) {
		ink = 'r'
	}
	t := h.text(fmt.Sprintf("Dia %d  %02d:%02d", c.Day+1, int(c.Hour), int(c.Hour*60)%60), ink)
	y := 8 + (rowH+2)*scale
	cv.Blit(8, y, h.panel(t.Width/scale+12, rowH, false))
	cv.Blit(8+6*scale, y+(rowH*scale-t.Height)/2, t)
}
