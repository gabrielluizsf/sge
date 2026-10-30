// Package bake renders 3D models into the sprites of an isometric game,
// ahead of time: a glTF model, posed at a moment of one of its
// animations, seen from the game's own corner (an iso.Projection), lit by
// a sun that casts shadows on it and on the ground under it, smoothed by
// supersampling. Games stay 2D and fast; their pictures come from 3D.
package bake

import "github.com/gabrielluizsf/sge/iso"

// Options are how a model is baked.
type Options struct {
	// Projection is the game's: a baked sprite drops onto its tiles.
	Projection iso.Projection
	// Unit is how many of the game world's units a model's unit is: the
	// side of a cell, for a model made in metres on a one-metre grid.
	Unit float64
	// Sun is the way toward the sun in the game's world, X along a row,
	// Y down a column, Z up; it need not be of length 1.
	Sun [3]float64
	// Light is the sun's strength, Sky the light from all round above,
	// and Ground that bounced from below.
	Light, Sky, Ground [3]float64
	// Exposure scales the light before it is tone mapped, as a camera's
	// does: at its default, ground facing the sun shows its texture's own
	// colour.
	Exposure float64
	// Shadow is how dark the shadow cast on the ground is, 0 for none.
	Shadow float64
	// Samples is how many samples a pixel takes across and down.
	Samples int
	// Water, if set, makes the ground's plane water's surface.
	Water *Water
}

// DefaultOptions are a morning sun from the upper left, as farm games
// light their pictures, a soft blue sky, and smooth edges.
func DefaultOptions(p iso.Projection) Options {
	return Options{
		Projection: p, Unit: p.Cell,
		Sun:   [3]float64{-0.35, 0.6, 1},
		Light: [3]float64{1.3, 1.25, 1.18}, Sky: [3]float64{0.4, 0.46, 0.6}, Ground: [3]float64{0.22, 0.2, 0.16},
		Exposure: 0.75, Shadow: 0.38, Samples: 3,
	}
}
