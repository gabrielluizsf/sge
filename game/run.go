package game

import "github.com/gabrielluizsf/antui"

// Run shows the world in a window of that title and size, stepping and
// drawing it every frame until the window is closed or Quit is called.
func (g *Game) Run(title string, width, height int) error {
	win, err := antui.Open(title, width, height)
	if err != nil {
		return err
	}
	g.window = win
	defer func() { g.window = nil }()
	for !g.quit && win.Begin() {
		dt := win.Delta()
		if g.FixedStep > 0 {
			dt = g.FixedStep
		}
		g.Step(dt)
		g.Draw(win.Canvas())
		win.End()
	}
	win.Close()
	return nil
}
