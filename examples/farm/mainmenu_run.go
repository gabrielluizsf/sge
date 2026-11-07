package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/game"
)

// attract makes s an attract mode with the main menu of the game titled
// title over it: the world plays on, its farmer the bot's, and the menus
// are the player's. It is the menu, to be driven.
func (s *Session) attract(title string) *mainMenu {
	m := newMainMenu(s.f, s.Game, title)
	s.f.ui.attract = true
	o := game.NewObject("mainmenu")
	o.Z, o.Screen, o.Always = zHUD+1, true, true
	o.AddScript(m)
	s.Game.Add(o)
	return m
}

// Menu shows the main menu of the game titled title in the open window
// win, over a film of the farm drawn by a, until the player
// chooses; closing the window or "Sair" is StartQuit. The window stays
// open for the game chosen.
func Menu(win *antui.Window, a Art, title string) Pick {
	s := New(a, Cinema())
	m := s.attract(title)
	if !s.Show(win) {
		return Pick{Start: StartQuit}
	}
	return m.pick
}
