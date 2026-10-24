package tilemap

// Animation is a tile that shows other tiles of its set in turn, FPS of
// them a second. Staggered starts each cell somewhere else in the turn, so
// a field of grass sways rather than marching.
type Animation struct {
	Frames    []int
	FPS       float64
	Staggered bool
}
