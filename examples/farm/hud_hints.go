package farm

import "github.com/gabrielluizsf/antui/canvas"

// hints tells the keys: what E does where the player stands, over their
// head, and to walk and use things at the top right.
func (h *hud) hints(cv *canvas.Canvas) {
	if words, _, ok := h.f.interaction(h.farmer); ok && words != "" {
		t := h.text(words, 'y')
		// Over the farmer's head: their position is where they stand.
		at := h.farmer.Game().ToScreen(h.farmer.Transform.Position)
		cv.Blit(int(at.X)-t.Width/2, int(at.Y)-tile-t.Height-scale*6, t)
	}
	keys := "Setas: andar   1-4: item   Espaço/clique: usar   E: interagir"
	if h.f.room != nil {
		keys = "Saia pelo tapete da porta"
	}
	t := h.text(keys, 'w')
	cv.Blit(cv.Width-t.Width-8, 8, t)
}
