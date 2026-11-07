package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/vignette"
)

// hurtRed is the red of a blow on the edges of the screen.
var hurtRed = vignette.Config{Color: 0xFFC01820, Alpha: 150, From: 0.68}

// hurtEdge reddens the edges of the screen after a blow, strongest at
// first and fading as the flash runs out.
func (h *hud) hurtEdge(cv *canvas.Canvas) {
	if h.f.flash <= 0 {
		return
	}
	if h.edge == nil || h.edge.Width != cv.Width || h.edge.Height != cv.Height {
		h.edge = vignette.New(cv.Width, cv.Height, hurtRed)
	}
	h.edge.Paint(cv, h.f.flash/flashFor)
}
