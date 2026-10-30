package farm

import "github.com/gabrielluizsf/antui/canvas"

// age is who the player is and how old, under the time.
func (h *hud) age(cv *canvas.Canvas) {
	f := h.f
	t := h.text(f.self.name+", "+itoa(f.self.who(f.clock.Day).Years)+" anos", 'k')
	y := 8 + 2*(rowH+2)*scale
	cv.Blit(8, y, h.panel(t.Width/scale+12, rowH, false))
	cv.Blit(8+6*scale, y+(rowH*scale-t.Height)/2, t)
}
