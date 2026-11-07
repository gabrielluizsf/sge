package main

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

// modelOf is the building the sprite name stands for: the stall, or a
// house by its kind.
func modelOf(name string) building.Building {
	switch name {
	case "stall":
		return building.Stall
	case "well":
		return building.Well
	case "stable":
		return building.Stable
	case "bakery":
		return building.Bakery
	case "butcher":
		return building.Butcher
	}
	var kind int
	fmt.Sscanf(name, "house%d", &kind)
	return building.Houses[kind%len(building.Houses)]
}

// Footprint is the cells the baked building name stands on: the metres its
// model covers, porch and steps with it, rounded up, a cell being a
// metre. Its walls are at the back of them, a porch in front.
func (a *realArt) Footprint(name string) (cols, rows int) {
	b := modelOf(name)
	return int(math.Ceil(b.Width)), int(math.Ceil(b.Depth + b.Front))
}

// pivotShift is how many pixels, across and down, a model's picture is
// moved so that the middle of its walls, where the baker puts its pivot,
// stands where it does in its footprint: a porch's cells in front of the
// walls put their middle, in the world, nearer the viewer than the
// footprint's, or farther when the model is turned to face away.
func (a *realArt) pivotShift(name string, heading int) (dx, dy float64) {
	_, rows := a.Footprint(name)
	s := (float64(rows) - modelOf(name).Depth) / 2
	if heading == backHeading {
		s = -s
	}
	return -kit.Cell * s, kit.Cell / 2 * s
}

// treeReach is how far beyond a cell's picture a tree grown reaches: a
// crown twenty metres across and as tall, in pixels, with its shadow.
const treeReach = 480

// Reach is treeReach: the trees are baked at their real size.
func (a *realArt) Reach() int { return treeReach }
