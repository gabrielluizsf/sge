package main

import (
	"time"

	"github.com/gabrielluizsf/sge/game"
)

// previewWait is how long the character must stay as it is before the
// person is made again; a change of colours only does not wait.
const previewWait = 60 * time.Millisecond

// creatorZ puts the screen over everything of the menu's.
const creatorZ = 1 << 20

// NewGame shows the character creation screen over the main menu's world,
// which plays on behind it, and calls done with the kit.Character the
// player makes. The person is made in the background, one worker, and
// the other bakes wait meanwhile.
func (a *realArt) NewGame(g *game.Game, done func(any)) { a.screenFor(g, done) }

// screenFor adds the creation screen to g and returns it.
func (a *realArt) screenFor(g *game.Game, done func(any)) *creator {
	var live *livePreview
	var release func()
	if w := a.world; w != nil {
		w.baker.Hold.Store(true)
		release = func() { w.baker.Hold.Store(false) }
		live = newLivePreview(bakerMaker{w.baker}, previewWait)
	}
	c := newCreator(a, live, done)
	c.release = release
	o := game.NewObject("creator")
	o.Screen, o.Always, o.Z = true, true, creatorZ
	o.AddScript(c)
	g.Add(o)
	return c
}
