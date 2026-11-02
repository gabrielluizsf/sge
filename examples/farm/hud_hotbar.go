package farm

import "github.com/gabrielluizsf/antui/canvas"

// hotbar is the bar at the bottom: what the player has chosen to have at
// hand, each numbered by its key, the one held in gold; with the
// inventory open, the whole screen of what is carried over it. It is out
// of the way while someone talks, whose box takes its place.
func (h *hud) hotbar(cv *canvas.Canvas) {
	if h.f.talk != nil {
		return
	}
	h.f.pack.Paint(cv, h.f.art)
}
