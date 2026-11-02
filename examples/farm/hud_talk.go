package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/inventory"
)

// talk is the conversation, role-playing style: the box at the bottom of
// the screen with the villager's portrait on its left, what they say in
// it and what the player may answer as buttons on its right; over its top
// right corner, how fond of the player they are, in hearts, and how far
// the bond has come in words.
func (h *hud) talk(cv *canvas.Canvas) {
	f := h.f
	v := f.talk.v
	b := f.talkBox()
	b.Paint(cv, f.art, f.art.Portrait(v.who))
	fr := b.Frame()
	label := h.text(heartsLabel(v), 'k')
	pw := label.Width + (heartsW+24)*scale
	x, y := fr.X+fr.W-pw, fr.Y-(rowH+2)*scale
	cv.Blit(x, y, h.panel(pw/scale, rowH, false))
	cv.Blit(x+6*scale, y+(rowH*scale-label.Height)/2, label)
	h.hearts(cv, v, [2]int{x + label.Width + 12*scale - heartsAt*scale, y})
}

// wants are the lines of a wish: each item, how many, and how many the
// player has.
func wants(f *farm, l inventory.List) []string {
	var out []string
	for _, s := range l {
		out = append(out, fmt.Sprintf("   %d %s  (você tem %d)", s.N, itemLabels[s.Item], f.bag.Count(s.Item)))
	}
	return out
}
