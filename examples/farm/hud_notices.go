package farm

import "github.com/gabrielluizsf/antui/canvas"

// notices are the news, in the middle of the screen, the latest lowest.
func (h *hud) notices(cv *canvas.Canvas) {
	for i, n := range h.f.notices {
		t := h.text(n.text, 'k')
		w := t.Width/scale + 16
		x, y := (cv.Width-w*scale)/2, cv.Height/3+i*(rowH+4)*scale
		cv.Blit(x, y, h.panel(w, rowH+2, true))
		cv.Blit(x+8*scale, y+((rowH+2)*scale-t.Height)/2, t)
	}
}
