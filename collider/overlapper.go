package collider

// Overlapper is a shape that decides itself which shapes it overlaps,
// rather than being tested by its bounds: a tile map's solid tiles, say.
type Overlapper interface {
	Collider
	Overlaps(other Collider) bool
}
