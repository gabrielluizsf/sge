package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// newStall is the shop's stall, standing on its two cells, stallCol and
// the one beside it on stallRow, among the farmer and the trees.
func newStall() *game.Object {
	src := spriteSVG("stall")
	g := art.Stall()
	w, h := len([]rune(g[0])), len(g)
	o := game.NewObject("stall")
	o.Z = zThings
	o.Sprite = &sprite.Sprite{Image: picture(src, w, h), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	o.Transform.Position = vector.Vector3{X: float64(stallCol+1) * tile, Y: (float64(stallRow) + 0.5) * tile}
	return o
}
