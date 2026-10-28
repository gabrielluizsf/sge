package main

import (
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// feetBox is the farmer's collider: their feet, around the point they
// stand on, their position, so they may walk right up to what stands in
// front of them or behind.
var feetBox = collider.Box{X: -6 * scale / 2, Y: -4 * scale / 2, Width: 6 * scale, Height: 4 * scale}

// feet is the point the farmer stands on.
func feet(o *game.Object) vector.Vector3 {
	return o.Transform.Position
}
