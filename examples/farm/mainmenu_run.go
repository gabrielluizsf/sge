package farm

import "github.com/gabrielluizsf/sge/game"

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

// Menu shows the main menu of the game titled title in a window w by h
// pixels, over the farm drawn by a playing by itself, until the player
// chooses; closing the window or "Sair" is StartQuit.
func Menu(a Art, title string, w, h int) (Pick, error) {
	s := New(a, Bot(10))
	m := s.attract(title)
	if err := s.Run(title, w, h); err != nil {
		return Pick{}, err
	}
	return m.pick, nil
}
