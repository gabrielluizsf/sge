package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// hotbar is the bar at the bottom: the seeds in hand, the can, with the
// water it holds, the axe and the hoe, the one held in yellow and each
// numbered by its key, then the wood and the harvest carried. With no
// axe yet, its slot is a question.
func (h *hud) hotbar(cv *canvas.Canvas) {
	const size, gap = 20 * scale, 22 * scale
	f := h.f
	carried := h.carried()
	n := slots + len(carried)
	x0 := (cv.Width - (n*gap - (gap - size))) / 2
	y := cv.Height - size - 8
	for i := range n {
		x := x0 + i*gap
		cv.Blit(x, y, h.panel(20, 20, i == f.slot))
		icon, count := "", -1
		switch {
		case i == slotSeeds:
			icon, count = "icon_"+crops[f.seed].seedIcon, f.bag.Count(seedsOf(f.seed))
		case i == slotCan:
			icon, count = "icon_can", f.inCan()
		case i == slotHoe:
			icon = "icon_hoe"
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

// carried are the items shown after the tools: the wood, and every crop
// picked that is in the bag.
func (h *hud) carried() []string {
	items := []string{wood}
	for c := range crops {
		if h.f.bag.Count(cropOf(c)) > 0 {
			items = append(items, cropOf(c))
		}
	}
	return items
}
