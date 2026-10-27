// Package art is the farm's pixel art: a fixed palette, and every sprite
// as a grid of its letters, written out as SVG. The farm is seen from a
// corner: its tiles are 32 by 48 pictures of a 16 by 16 cell's diamond,
// ground projected flat on it and what stands on it drawn upright.
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
	'h': "#2a5a33", // shadow on grass
	// The farmer's: straw, leather, hair, shirt, denim and brass.
	'g': "#e2bd6e", // straw
	'G': "#b88d45", // straw, shade
	'z': "#7a4a2b", // leather
	'Z': "#4d2c18", // leather, dark
	'l': "#6b3e22", // hair
	'L': "#43260f", // hair, dark
	'f': "#f1e6d0", // shirt
	'F': "#cdbd9f", // shirt, shade
	'j': "#3e65a6", // denim
	'J': "#2a4677", // denim, dark
	'i': "#6a93cf", // denim, light
	'B': "#d6a644", // brass
}

// Size is the side of every sprite, in pixels.
const Size = 16
