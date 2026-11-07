package farm

import "github.com/gabrielluizsf/sge/game"

// NewGameArt is an Art with a step before a new game: a character
// creation screen, say. The main menu's "Novo jogo" calls NewGame with the
// game the menu runs in, whose window it may add objects to (set Always on
// them, for their scripts run while the menu holds the world), and waits
// for done: the screen removes its objects and calls done with what the
// player chose, or nil, and the new game starts. An Art without it starts
// the new game at once.
type NewGameArt interface {
	NewGame(g *game.Game, done func(choice any))
}

// ChosenArt is an Art that wants to know what the player chose in
// NewGame, before the new game is built: Main calls Chosen with it.
type ChosenArt interface {
	Chosen(choice any)
}
