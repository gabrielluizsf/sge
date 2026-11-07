package bakes

import (
	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/farm/kit"
)

// Patches are the cells a ground's texture spans: two by two metres, so
// its four tiles laid by the cells' parity meet without a seam.
var Patches = []bake.Patch{{Span: 2, Col: 0, Row: 0}, {Span: 2, Col: 1, Row: 0}, {Span: 2, Col: 0, Row: 1}, {Span: 2, Col: 1, Row: 1}}

// Headings are the turns each model is baked at.
var Headings = []float64{45, 135, 225, 315}

// Options are the farm's light: a morning sun from the upper left.
func Options() bake.Options {
	o := bake.DefaultOptions(kit.Projection)
	o.Unit = kit.Cell
	return o
}
