package farm

import "github.com/gabrielluizsf/antui/canvas"

// text is s in the HUD's face, in the palette's ink, as the art writes
// it.
func (h *hud) text(s string, ink rune) *canvas.Canvas { return h.f.art.Text(s, ink) }

// panel is a HUD panel of w by ht pixels of art, gold when selected.
func (h *hud) panel(w, ht int, selected bool) *canvas.Canvas {
	return h.f.art.Panel(w, ht, selected)
}

// icon is the HUD's icon name.
func (h *hud) icon(name string) *canvas.Canvas { return h.f.art.Icon(name) }
