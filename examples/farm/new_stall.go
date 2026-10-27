package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// newStall is the shop's stall, its art drawn at half its size so that
// it stands on one cell, stallCol, stallRow, among the farmer and the
// trees.
func newStall() *game.Object {
	o := game.NewObject("stall")
	o.Z = zThings
	o.Sprite = &sprite.Sprite{Image: picture(spriteSVG("stall"), art.StallW/2, art.StallH/2), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	o.Transform.Position = vector.Vector3{X: (float64(stallCol) + 0.5) * tile, Y: (float64(stallRow) + 0.5) * tile}
	return o
}
