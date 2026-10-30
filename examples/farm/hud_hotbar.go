package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// hotbar is the bar at the bottom: the two seeds, the can and the axe,
// the one held in yellow and each numbered by its key, then the harvest
// and the wood carried. With no axe yet, its slot is a question.
func (h *hud) hotbar(cv *canvas.Canvas) {
	const size, gap = 20 * scale, 22 * scale
	f := h.f
	carried := append([]string{}, cropOf(0), cropOf(1), wood)
	n := slots + len(carried)
	x0 := (cv.Width - (n*gap - (gap - size))) / 2
	y := cv.Height - size - 8
	for i := range n {
		x := x0 + i*gap
		cv.Blit(x, y, h.panel(20, 20, i == f.slot))
		icon, count := "", -1
		switch {
		case i < slotCan:
			icon, count = "icon_"+crops[i].seedIcon, f.bag.Count(seedsOf(i))
		case i == slotCan:
			icon = "icon_can"
		case i == slotAxe && f.axe == noAxe:
			q := h.text("?", 'k')
			cv.Blit(x+(size-q.Width)/2, y+(size-q.Height)/2, q)
		case i == slotAxe:
			icon = "icon_axe_" + axes[f.axe].name
		default:
			item := carried[i-slots]
			icon, count = "icon_"+item, f.bag.Count(item)
		}
		if icon != "" {
			cv.Blit(x+2*scale, y+2*scale, h.icon(icon))
		}
		if i < slots {
			k := h.text(fmt.Sprint(i+1), 'w')
			cv.Blit(x+(size-k.Width)/2, y-k.Height, k)
		}
		if count >= 0 {
			t := h.text(fmt.Sprint(count), 'w')
			cv.Blit(x+size-t.Width+outline, y+size-t.Height+2*outline, t)
		}
	}
}
