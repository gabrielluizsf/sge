package farm

import "github.com/gabrielluizsf/antui"

// shopKeys runs the open shop's menu: up and down choose a line, space or
// enter does it, and E or escape closes the shop.
func (f *farm) shopKeys(win *antui.Window) {
	s := f.shop
	lines := sellChoice + 1
	if s.animals {
		lines = len(stableOffers)
	}
	switch {
	case win.KeyPressed(antui.KeyEscape) || win.KeyPressed(antui.KeyE):
		s.open = false
		f.sound("ui_close", nil)
	case win.KeyPressed(antui.KeyUp) || win.KeyPressed(antui.KeyW):
		s.choice = (s.choice + lines - 1) % lines
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeyDown) || win.KeyPressed(antui.KeyS):
		s.choice = (s.choice + 1) % lines
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeySpace) || win.KeyPressed(antui.KeyEnter):
		f.choose(f.kit, s)
	}
}
