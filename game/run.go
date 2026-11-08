package game

import "github.com/gabrielluizsf/antui"

// Run shows the world in a window of that title and size, stepping and
// drawing it every frame until the window is closed or Quit is called,
// and then closes the window. To show one game after another in the same
// window, open it with antui.Open and Show each in it.
func (g *Game) Run(title string, width, height int) error {
	win, err := antui.Open(title, width, height)
	if err != nil {
		return err
	}
	g.Show(win)
	win.Close()
	return nil
}
