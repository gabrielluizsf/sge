package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// hotbar is the bar at the bottom: the two seeds and the can, the one held
// in yellow and each numbered by its key, then the harvest carried.
func (h *hud) hotbar(cv *canvas.Canvas) {
	const size, gap = 20 * scale, 22 * scale
	n := slots + len(crops)
	x0 := (cv.Width - (n*gap - (gap - size))) / 2
	y := cv.Height - size - 8
	for i := range n {
		x := x0 + i*gap
		cv.Blit(x, y, panel(20, 20, i == h.f.slot))
		var icon string
		count := -1
		switch {
		case i < slotCan:
			icon, count = "icon_"+crops[i].seedIcon, h.f.seeds[i]
		case i == slotCan:
			icon = "icon_can"
		default:
			c := i - slots
			icon, count = "icon_"+crops[c].cropIcon, h.f.harvest[c]
		}
		cv.Blit(x+2*scale, y+2*scale, look(icon))
		if i < slots {
			cv.Blit(x+size/2-3*scale, y-10*scale, text(fmt.Sprint(i+1), 'w'))
		}
		if count >= 0 {
			t := text(fmt.Sprint(count), 'w')
			cv.Blit(x+size-t.Width-scale, y+size-t.Height, t)
		}
	}
}
