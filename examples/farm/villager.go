package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/request"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// villager is one who lives in a town, the smith or not: who they are,
// what they have, their home and what they ask of the player. They live
// on while their town sleeps; their object is in the game only while it
// is awake.
type villager struct {
	name  string
	woman bool
	role  string
	look  art.Look
	kit   *kit
	town  *town.Town
	home  *town.House // nil for a smith with no house in town
	order request.Order
	wait  float64 // seconds until they ask for something
	pos   vector.Vector3
	obj   *game.Object
	mind  *mind
	luck  *random.Source
	// inside is whether they are in their house, for homeFor more seconds
	// at least, and all night when they have gone to bed.
	inside  bool
	homeFor float64
	// frames are their pictures' SVG by name, drawn when first seen.
	frames map[string]string
	// fuel is how much wood a smith's forge has left to burn.
	fuel float64
}

// smith reports whether v is their town's smith.
func (v *villager) smith() bool { return v.role == roleSmith }

// wants reports whether v has something to ask of the player.
func (v *villager) wants(f *farm) bool {
	if v.smith() {
		return !f.errandDone
	}
	return len(v.order.Want) > 0
}

// ask gives v a new order once they have waited long enough since the
// last.
func (v *villager) ask(dt float64) {
	if v.smith() || len(v.order.Want) > 0 {
		return
	}
	if v.wait -= dt; v.wait <= 0 {
		v.order = orderBook.Draw(v.luck)
	}
}
