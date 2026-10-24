package main

import (
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// feetBox is the farmer's collider: their feet, so their head may pass in
// front of what is behind them.
var feetBox = collider.Box{X: -6 * scale / 2, Y: 4 * scale, Width: 6 * scale, Height: 4 * scale}

// feet is the point the farmer stands on.
func feet(o *game.Object) vector.Vector3 {
	return o.Transform.Position.Add(vector.Vector3{Y: feetBox.Y + feetBox.Height/2})
}
