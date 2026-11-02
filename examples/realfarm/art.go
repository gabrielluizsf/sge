package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// realArt is the farm drawn from baked sprites: the ground, water, crops,
// trees, houses and people read from the shelf. What the shelf does not
// have, the HUD's text, panels and icons and the insides of the houses,
// is the stylized art's, but the text and panels, which are Source Serif on dark glass.
type realArt struct {
	farm.Art // the stylized art, for what is not baked
	shelf    *bakes.Shelf
	seen     map[string]*baked
	hud      map[string]*canvas.Canvas // the HUD's text and panels, by what they draw
}

// newArt is the realistic look, reading sprites from shelf as the game
// asks for them, which the baker may still be adding to.
func newArt(shelf *bakes.Shelf) farm.Art {
	return &realArt{Art: farm.SVGArt(), shelf: shelf, seen: map[string]*baked{}, hud: map[string]*canvas.Canvas{}}
}

// Gridded is true: the ground is laid by the cells' parity, its patches
// meeting without a seam.
func (a *realArt) Gridded() bool { return true }

// Houses is how many kinds of house the baker builds.
func (a *realArt) Houses() int { return len(building.Houses) }
