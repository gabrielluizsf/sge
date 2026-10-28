package main

import "github.com/gabrielluizsf/antui"

// talkKeys runs the open conversation's menu: up and down choose a line,
// space or enter says it, and E or escape ends the conversation.
func (f *farm) talkKeys(win *antui.Window) {
	t := f.talk
	ls := f.lines(t.v)
	t.choice = min(t.choice, len(ls)-1)
	switch {
	case win.KeyPressed(antui.KeyEscape) || win.KeyPressed(antui.KeyE):
		f.talk = nil
		f.sound("ui_close", nil)
	case win.KeyPressed(antui.KeyUp) || win.KeyPressed(antui.KeyW):
		t.choice = (t.choice + len(ls) - 1) % len(ls)
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeyDown) || win.KeyPressed(antui.KeyS):
		t.choice = (t.choice + 1) % len(ls)
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeySpace) || win.KeyPressed(antui.KeyEnter):
		ls[t.choice].do(f)
	}
}
