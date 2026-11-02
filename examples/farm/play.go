package farm

// play shows the main menu and then plays what it chose, in a window w
// by h pixels. It exits on an error.
func play(a Art, title string, w, h int) {
	p, err := Menu(a, title, w, h)
	exitOn(err)
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
	exitOn(s.Run(title, w, h))
}
