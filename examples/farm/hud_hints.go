package main

import "github.com/gabrielluizsf/antui/canvas"

// hints tells the keys: what E does where the player stands, over their
// head, and to walk and use things at the top right.
func (h *hud) hints(cv *canvas.Canvas) {
	if words, _, ok := h.f.interaction(h.farmer); ok && words != "" {
		t := text(words, 'y')
		// Over the farmer's head: their position is where they stand.
		at := h.farmer.Game().ToScreen(h.farmer.Transform.Position)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-tile-t.Height-scale*6, t)
	}
	keys := "SETAS: ANDAR  1-4: ITEM  ESPACO/CLIQUE: USAR  E: FALAR"
	if h.f.room != nil {
		keys = "SAIA PELO TAPETE DA PORTA"
	}
	t := text(keys, 'w')
	cv.Blit(cv.Width-t.Width-8, 8, t)
}
