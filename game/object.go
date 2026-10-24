package game

import (
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/transform"
)

// Object is one thing in the world.
type Object struct {
	Name      string
	Transform transform.Transform
	// Sprite is what the object looks like, unless Animator has a frame
	// to show instead.
	Sprite   *sprite.Sprite
	Animator *animation.Animator
	// Collider is the object's shape, relative to its position; an object
	// without one touches nothing.
	Collider collider.Collider
	// Z orders drawing: lower first, so higher is on top.
	Z      int
	Hidden bool
	// Screen fixes the object to the screen: its position is a point of
	// the canvas, where the camera does not move it, as a HUD's is.
	Screen bool

	scripts   []Script
	handles   bool // a script wants to be told what the object touches
	game      *Game
	started   bool
	destroyed bool

	// world is WorldCollider's last answer, for worldFrom placed at
	// worldX, worldY: collisions ask for it once per pair.
	world          collider.Collider
	worldFrom      collider.Collider
	worldX, worldY float64
}
