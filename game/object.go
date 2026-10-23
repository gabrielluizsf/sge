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

	scripts   []Script
	game      *Game
	started   bool
	destroyed bool
}
