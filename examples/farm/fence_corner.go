package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// cornerHeadings are the heading each corner of a fence is baked at: the
// model's rails run from its post along its +X and its +Z.
var cornerHeadings = map[string]int{"fence_ne": 0, "fence_se": 1, "fence_sw": 2, "fence_nw": 3}

// fenceCorner is the corner name of a fence, or, while it is not baked, the
// straight run, so a fence is never seen with a hole in it.
func (a *realArt) fenceCorner(name string) *canvas.Canvas {
	if b, ok := a.bake(fmt.Sprintf("building_fence_corner_%d", cornerHeadings[name])); ok {
		return b.standing()
	}
	return a.buildingTile("fence_straight")
}
