package main

// scale is how many screen pixels a pixel of the art takes, tile the side
// of a cell in the world, and isoW by isoH a tile's picture on screen: a
// cell's diamond, twice as wide as it is tall, and room above it for
// what stands on it.
const (
	scale = 2
	tile  = 16 * scale
	isoW  = 2 * tile
	isoH  = 3 * tile
)
