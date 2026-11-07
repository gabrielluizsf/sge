package main

import (
	"sync/atomic"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
)

// world is the baker and the bakes in the background behind the game.
type world struct {
	baker *bakes.Baker
	stop  *atomic.Bool  // set to stop the background bakes under way
	win   *antui.Window // the game's one window, nil in a demo, which has none to load in
}

// start bakes what the game cannot begin without, in a loading screen,
// and the rest in the background.
func (w *world) start() error {
	core, rest := plan(w.baker)
	if w.win != nil {
		if err := boot(w.win, w.baker, core); err != nil {
			return err
		}
		w.stop = new(atomic.Bool)
		go background(w.baker, rest, w.stop)
	}
	return nil
}

// restart stops the background bakes and starts over from what the game
// needs now: the player having a new look.
func (w *world) restart() error {
	if w.stop != nil {
		w.stop.Store(true)
	}
	return w.start()
}
