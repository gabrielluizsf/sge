package farm

import "github.com/gabrielluizsf/antui"

// questionKeys runs the open question: up and down, or left and right,
// choose an answer, and space or enter gives it.
func (f *farm) questionKeys(win *antui.Window) {
	q := f.question
	n := len(q.answers)
	switch {
	case win.KeyPressed(antui.KeyUp) || win.KeyPressed(antui.KeyW) || win.KeyPressed(antui.KeyLeft) || win.KeyPressed(antui.KeyA):
		q.choice = (q.choice + n - 1) % n
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeyDown) || win.KeyPressed(antui.KeyS) || win.KeyPressed(antui.KeyRight) || win.KeyPressed(antui.KeyD):
		q.choice = (q.choice + 1) % n
		f.sound("ui_tick", nil)
	case win.KeyPressed(antui.KeySpace) || win.KeyPressed(antui.KeyEnter):
		f.answerQuestion()
	}
}
