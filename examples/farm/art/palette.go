// Package art is the farm's pixel art: a fixed palette, and every sprite
// as a 16 by 16 grid of its letters, written out as SVG.
package art

// Palette is every colour the farm is drawn in, by the letter its grids
// use. A sprite with any other colour is not the farm's: palette_test.go
// says so. '.' is transparent.
var Palette = map[rune]string{
	'k': "#2b1f2e", // outline
	'1': "#2f6b3a", // grass, dark
	'2': "#4e9a4a", // grass
	'3': "#8bc65a", // grass, light
	'a': "#6b4226", // earth, dark
	'b': "#9c6b3c", // earth
	'c': "#c99a5b", // earth, light
	'e': "#4a2c1a", // wet soil
	'u': "#2a4d8f", // water, deep
	'v': "#3f7fd1", // water
	'x': "#8fd3ff", // water, light
	'o': "#7a4b2a", // wood, dark
	'q': "#b07a45", // wood
	's': "#f2c39b", // skin
	't': "#c98b63", // skin, shade
	'r': "#c0392b", // red
	'R': "#8e2920", // red, dark
	'n': "#f28c28", // orange
	'N': "#c45f14", // orange, dark
	'y': "#f4d35e", // yellow
	'Y': "#c9a227", // yellow, dark
	'w': "#f4f1e8", // white
	'm': "#7d7d7d", // metal, dark
	'M': "#b8b8b8", // metal
	'p': "#e86a92", // pink
}

// Size is the side of every sprite, in pixels.
const Size = 16
