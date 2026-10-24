package main

import "github.com/gabrielluizsf/antui/canvas"

// hints tells the keys: to open the shop at the stall, and to walk and
// use things anywhere else.
func (h *hud) hints(cv *canvas.Canvas) {
	if h.f.shop.open {
		return
	}
	if h.f.atStall(h.farmer) {
		t := text("E: LOJA", 'y')
		at := h.farmer.Game().ToScreen(h.farmer.Transform.Position)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-tile, t)
	}
	t := text("SETAS: ANDAR  1-3: ITEM  ESPACO: USAR", 'w')
	cv.Blit(8, cv.Height-t.Height-8-24*scale, t)
}
