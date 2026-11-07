package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/roundpic"
)

// The portrait's ring and the grey of a face not yet baked.
const (
	ringInk  canvas.Color = 0xFF5A3A22
	blankInk canvas.Color = 0xFFA9ADB3
)

// portrait is the player's face cut to a circle d pixels across, with its
// ring. It is made when the art's picture or the size changes, not every
// frame; a face not yet baked shows a neutral circle that is replaced
// when it is.
func (h *hud) portrait(d int) *canvas.Canvas {
	src := h.f.art.Portrait(h.f.self.who(h.f.clock.Day))
	if h.round != nil && h.roundSrc == src && h.roundD == d {
		return h.round
	}
	if src == nil {
		h.round = roundpic.Neutral(d, blankInk, ringInk)
	} else {
		h.round = roundpic.Crop(src, d, ringInk)
	}
	h.roundSrc, h.roundD = src, d
	return h.round
}
