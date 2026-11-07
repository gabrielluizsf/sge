package farm

import "github.com/gabrielluizsf/antui"

// play opens the game's one window, titled title, w by h pixels, and plays
// in it: the main menu, then what it chose. It exits on an error.
func play(a Art, title string, w, h int) {
	win, err := antui.Open(title, w, h)
	exitOn(err)
	Play(win, a, title)
	win.Close()
}

// Play shows the main menu of the game titled title in the open window
// win, and then plays what it chose in the same window: scenes change, the
// window stays. It returns when the player quits or closes the window.
func Play(win *antui.Window, a Art, title string) {
	p := Menu(win, a, title)
	if p.Start == StartQuit {
		return
	}
	if c, ok := a.(ChosenArt); ok {
		c.Chosen(p.Choice)
	}
	s := New(a, Keyboard())
	if p.Start == StartContinue {
		exitOn(s.Load(title))
	}
	s.EnableSave(title)
	s.Show(win)
}
