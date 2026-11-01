package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/sprite"
)

// The headings a house is baked at that show its door toward the viewer
// and away from them.
const (
	frontHeading = 0
	backHeading  = 2
)

// Building is the stall or the house called name, baked; a plain sprite
// of the stylized art's while the baker has not made it.
func (a *realArt) Building(name string) *sprite.Sprite {
	id, heading := buildingOf(name)
	b, ok := a.bake(fmt.Sprintf("building_%s_%d", id, heading))
	if !ok {
		return a.Art.Building(name)
	}
	return &sprite.Sprite{Image: b.cv, PivotX: float64(b.at.PivotX) / float64(b.at.W), PivotY: float64(b.at.PivotY) / float64(b.at.H)}
}

// buildingOf is the building and the heading of the sprite called name:
// "stall", or "house" and a kind, with "_back" after it when the door
// faces away.
func buildingOf(name string) (string, int) {
	if name == "stall" {
		return name, frontHeading
	}
	var kind int
	fmt.Sscanf(name, "house%d", &kind)
	id := building.Houses[kind%len(building.Houses)].Name
	if len(name) > 5 && name[len(name)-5:] == "_back" {
		return id, backHeading
	}
	return id, frontHeading
}
