// Package scene switches a game between scenes, the world outside and the
// inside of a house, say: a scene is the objects in the game while it
// plays, and leaving it keeps them as they are, to come back to.
package scene

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// Stage is a game's scenes. Keep, when set, reports the objects that are
// in every scene, the player and the HUD; Current is the scene playing.
type Stage struct {
	Game    *game.Game
	Keep    func(o *game.Object) bool
	Current string
	saved   map[string]*shot
}

// shot is a scene as it was left: its objects, and where the camera was
// and what was behind it.
type shot struct {
	objects    []*game.Object
	camera     vector.Vector3
	background canvas.Color
	covered    bool
}
