package main

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// Chosen is the character the player made in the creation screen: the
// player villager becomes them, they are kept for the games to come, and
// the player is baked in their look, the loading screen shown for the
// motions the game cannot begin without.
func (a *realArt) Chosen(choice any) {
	c, ok := choice.(kit.Character)
	if !ok || c.Check() != nil {
		return
	}
	kit.SetPlayer(c)
	if a.store != nil {
		if err := kit.SaveCharacter(a.store, c); err != nil {
			fmt.Fprintln(os.Stderr, "realfarm: the look was not saved:", err)
		}
	}
	if a.world != nil {
		if err := a.world.restart(); err != nil {
			fail(err)
		}
	}
}

// PlayerName is the name the player gave their character, which the
// farm calls the farmer by; "" while they gave none.
func (a *realArt) PlayerName() string { return kit.PlayerName() }
