package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// newStall is a shop's stall, its art drawn at half its size so that it
// stands on one cell, at, among the people and the trees.
func newStall(at town.Cell) *game.Object {
	o := game.NewObject("stall")
	o.Z = zThings
	o.Sprite = &sprite.Sprite{Image: picture(spriteSVG("stall"), art.StallW/2, art.StallH/2), PivotX: art.StallPivot[0], PivotY: art.StallPivot[1]}
	o.Transform.Position = vector.Vector3{X: (float64(at.Col) + 0.5) * tile, Y: (float64(at.Row) + 0.5) * tile}
	return o
}

// newHouse is a house of the given kind on its two by two cells from
// corner, its door toward the viewer when front is set.
func newHouse(kind int, corner town.Cell, front bool) *game.Object {
	o := game.NewObject("house")
	o.Z = zThings
	name := spriteName("house", kind, front)
	o.Sprite = &sprite.Sprite{Image: picture(spriteSVG(name), art.HouseW, art.HouseH), PivotX: art.HousePivot[0], PivotY: art.HousePivot[1]}
	o.Transform.Position = vector.Vector3{X: float64(corner.Col+1) * tile, Y: float64(corner.Row+1) * tile}
	return o
}

// spriteName is the name of a house's sprite of the kind, from the front
// or the back.
func spriteName(base string, kind int, front bool) string {
	s := base + itoa(kind)
	if !front {
		s += "_back"
	}
	return s
}
