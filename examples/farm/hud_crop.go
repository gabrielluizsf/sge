package farm

import "github.com/gabrielluizsf/antui/canvas"

// cropLabel writes over the world, above the planted plot the cursor
// marks, how its crop is doing.
func (h *hud) cropLabel(cv *canvas.Canvas) {
	p := h.f.markedPlot(h.farmer)
	if p == nil {
		return
	}
	t := h.text(cropNote(p), 'w')
	at := h.farmer.Game().ToScreen(h.f.ground.CellCenter(p.col, p.row))
	cv.Blit(int(at.X)-t.Width/2, int(at.Y)-tile-t.Height, t)
}
