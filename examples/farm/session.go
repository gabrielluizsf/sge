package farm

import "github.com/gabrielluizsf/sge/game"

// Session is a game of the farm: the engine's game, which plays and draws
// it, and the farm in it.
type Session struct {
	Game *game.Game
	f    *farm
}

// New is a game of the farm drawn by a, its farmer played by ctl: the
// farm laid out with the farmer in front of the stall, the towns around
// it, the clock, the night and the HUD. Its rules and its world are the
// same whatever the art.
func New(a Art, ctl Controller) *Session {
	g, f := build(a, ctl)
	return &Session{Game: g, f: f}
}

// Run plays s in a window titled title, w by h pixels, with its sounds,
// until the window is closed.
func (s *Session) Run(title string, w, h int) error {
	addSound(s.Game, s.f)
	return s.Game.Run(title, w, h)
}
