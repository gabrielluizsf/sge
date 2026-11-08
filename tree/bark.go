package tree

import (
	"image"
	"image/color"
)

// Bark is what a trunk is clad in: a PBR texture set, as Poly Haven's
// are, its colour, its glTF metal-rough map (occlusion in red, roughness
// in green) and its normal map, nil where it has none; how many metres
// one repeat of it covers; and its colour, in sRGB, times the texture's,
// or, with no texture, alone (zero for a plain brown).
type Bark struct {
	Albedo, ARM, Normal image.Image
	Span                float64
	Color               color.NRGBA
}

// Barks are the barks the species wear: a rough, furrowed pine's, a
// smooth one, and a palm's ringed one. A bark left empty is plain brown.
type Barks struct {
	Pine, Smooth, Palm Bark
}
