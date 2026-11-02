package main

import (
	"sync/atomic"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// world is the baker and the bakes in the background behind the game.
type world struct {
	baker *bakes.Baker
	stop  *atomic.Bool // set to stop the background bakes under way
	demo  bool         // a demo has no window to load in
}

// start bakes what the game cannot begin without, in a loading screen,
// and the rest in the background.
func (w *world) start() error {
	core, rest := plan(w.baker)
	if !w.demo {
		if err := boot(w.baker, core); err != nil {
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
