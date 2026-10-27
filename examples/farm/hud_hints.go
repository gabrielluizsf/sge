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
		// Over the farmer's head: their position is where they stand.
		at := h.farmer.Game().ToScreen(h.farmer.Transform.Position)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-tile-t.Height-scale*2, t)
	}
	// The keys at the top right, out of the way of the purse and the bar.
	t := text("SETAS: ANDAR  1-3: ITEM  ESPACO/CLIQUE: USAR", 'w')
	cv.Blit(cv.Width-t.Width-8, 8, t)
}
