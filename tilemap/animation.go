package tilemap

// Animation is a tile that shows other tiles of its set in turn, FPS of
// them a second. Staggered starts each cell somewhere else in the turn, so
// a field of grass sways rather than marching.
type Animation struct {
	Frames    []int
	FPS       float64
	Staggered bool
	// Parity shows each cell one frame for good, picked by where it lies:
	// Frames[col%2 + 2*(row%2)], the four patches of a texture that
	// spans two cells across, so they meet without a seam. FPS and
	// Staggered do not matter.
	Parity bool
}
