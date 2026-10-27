package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// hud draws what the farmer needs to see over the farm, fixed to the
// screen, as farm games lay it out: the money at the top left, what they can hold in a bar at
// the bottom, the cell they face marked, words rising from what happens,
// and the shop over it all when open.
type hud struct {
	f      *farm
	farmer *game.Object
	p      *farmer
}

func (h *hud) Update(*game.Object, float64) {}

func (h *hud) Paint(o *game.Object, cv *canvas.Canvas) {
	if !h.f.shop.open {
		h.cursor(cv)
	}
	for _, w := range h.f.floats {
		t := text(w.text, w.ink)
		at := o.Game().ToScreen(w.at)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-t.Height, t)
	}
	h.money(cv)
	h.hotbar(cv)
	h.hints(cv)
	if h.f.shop.open {
		h.shop(cv)
	}
}

// cursor marks the plot the farmer works: its diamond, outlined.
func (h *hud) cursor(cv *canvas.Canvas) {
	col, row := h.p.target(h.farmer)
	if h.f.plotAt(col, row) == nil {
		return
	}
	at := h.farmer.Game().ToScreen(h.f.ground.CellCenter(col, row))
	cv.Blit(int(at.X)-isoW/2, int(at.Y)-tile/2, picture(spriteSVG("cursor"), isoW/scale, tile/scale))
}

// money is the purse at the top left: a coin and what is in it.
func (h *hud) money(cv *canvas.Canvas) {
	cv.Blit(8, 8, panel(52, 14, false))
	cv.Blit(8-scale, 8-scale, look("icon_coin"))
	cv.Blit(8+16*scale, 8+3*scale, text(fmt.Sprint("$", h.f.money), 'w'))
}
