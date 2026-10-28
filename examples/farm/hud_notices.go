package main

import "github.com/gabrielluizsf/antui/canvas"

// notices are the news, in the middle of the screen, the latest lowest.
func (h *hud) notices(cv *canvas.Canvas) {
	for i, n := range h.f.notices {
		t := text(n.text, 'y')
		w := t.Width/scale + 12
		x, y := (cv.Width-w*scale)/2, cv.Height/3+i*14*scale
		cv.Blit(x, y, panel(w, 13, false))
		cv.Blit(x+6*scale, y+3*scale, t)
	}
}
