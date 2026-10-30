package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// bubbles are the bubbles over the heads of the villagers in sight who
// want something of the player, as cosy games show what is on one's mind.
func (h *hud) bubbles(cv *canvas.Canvas) {
	f := h.f
	show := func(v *villager, o *game.Object) {
		if !v.wants(f) || o.Hidden {
			return
		}
		at := h.farmer.Game().ToScreen(o.Transform.Position)
		b := h.icon("icon_bubble")
		cv.Blit(int(at.X)-b.Width/2, int(at.Y)-tile-b.Height-2*scale, b)
	}
	if f.room != nil {
		for _, r := range f.residents {
			show(r.v, r.obj)
		}
		return
	}
	for _, v := range f.awake() {
		show(v, v.obj)
	}
}
