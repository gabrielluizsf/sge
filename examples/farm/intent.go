package main

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// intent is what the farmer is told to do in a frame: walk towards dir,
// a way in the world, no further than upTo pixels when that is set; turn
// to face, when set; work the cell aim, when set, rather than the one in
// front; and use what they hold there.
type intent struct {
	dir  vector.Vector3
	upTo float64
	face string
	aim  *[2]int
	use  bool
}

// controller decides a person's intent each frame, dt seconds after the
// last: the keyboard, the demo's bot, or a villager's mind.
type controller interface {
	control(f *farm, o *game.Object, dt float64) intent
}
