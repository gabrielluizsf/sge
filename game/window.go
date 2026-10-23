package game

import "github.com/gabrielluizsf/antui"

// Window is the window the world is running in, for its input: nil when
// it is not running in one.
func (g *Game) Window() *antui.Window { return g.window }
