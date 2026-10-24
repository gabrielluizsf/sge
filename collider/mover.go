package collider

// Mover is a shape that knows how to be moved, which Offset asks it to be
// rather than moving its bounds.
type Mover interface {
	Collider
	Moved(dx, dy float64) Collider
}
