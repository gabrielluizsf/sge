package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/save"
)

// realArt is the farm drawn from baked sprites: the ground, water, crops,
// trees, houses, their insides and the people read from the shelf, and
// the HUD drawn as its own. Nothing of it is the stylized art's: what the
// baker has not made yet is blank.
type realArt struct {
	shelf *bakes.Shelf
	seen  map[string]*baked
	hud   map[string]*canvas.Canvas // the HUD's text and panels, by what they draw
	world *world                    // what bakes behind the game; nil with nothing to bake
	store *save.Store               // where the player's character is kept; nil to keep none
}

// newArt is the realistic look, reading sprites from shelf as the game
// asks for them, which the baker may still be adding to.
func newArt(shelf *bakes.Shelf) farm.Art {
	return &realArt{shelf: shelf, seen: map[string]*baked{}, hud: map[string]*canvas.Canvas{}}
}

// Gridded is true: the ground is laid by the cells' parity, its patches
// meeting without a seam.
func (a *realArt) Gridded() bool { return true }

// Houses is how many kinds of house the baker builds.
func (a *realArt) Houses() int { return len(building.Houses) }
