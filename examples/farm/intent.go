package main

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// intent is what the farmer is told to do in a frame: walk towards dir,
// no further than upTo pixels when that is set; turn to face, when set;
// and use what they hold on the cell they face.
type intent struct {
	dir  vector.Vector3
	upTo float64
	face string
	use  bool
}

// controller decides the farmer's intent each frame: the keyboard, or
// the demo's bot.
type controller interface {
	control(f *farm, o *game.Object) intent
}
