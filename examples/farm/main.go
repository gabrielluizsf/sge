// Command farm is the game: a farm drawn from real things: its ground, trees,
// crops, houses and people are baked from CC0 models into sprites when
// the game first starts (a few minutes, with a loading screen), kept in
// the user's cache and reused after. What is not needed to begin is baked
// while the game plays. A new game begins with the player making their
// character.
//
//	go run ./examples/farm
package main

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/i18n"
)

// gameTitle is the game's title, on its window and its menu; saveFolder is
// what names the folder of its saves, and stays as it first was so the saves
// made under that name go on being found.
const (
	gameTitle  = "sge farm"
	saveFolder = "sge realfarm"
)

func main() {
	b, err := newBaker()
	if err != nil {
		fail(err)
	}
	store, err := kit.CharacterStore(saveFolder)
	if err != nil {
		fail(err)
	}
	if c, ok := kit.LoadCharacter(store); ok { // a hostile or missing file leaves the default player
		kit.SetPlayer(c)
	}
	w := &world{baker: b}
	art := newArt(b.Shelf).(*realArt)
	art.world, art.store = w, store
	if os.Getenv("FARM_DEMO") != "" { // a demo plays what is baked, with no window
		farm.Main(art, saveFolder)
		return
	}
	i18n.Load()
	// One window for the whole game: the loading screen, the menu, the
	// character's making and the play follow one another in it.
	win, err := antui.Open(gameTitle, 1280, 720)
	if err != nil {
		fail(err)
	}
	w.win = win
	if err := w.start(); err != nil {
		fail(err)
	}
	farm.Play(win, art, saveFolder)
	win.Close()
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "farm:", err)
	os.Exit(1)
}
