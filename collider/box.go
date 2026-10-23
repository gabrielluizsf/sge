package collider

// Box is an axis-aligned rectangle with its top-left corner at X, Y.
type Box struct {
	X, Y, Width, Height float64
}

// Bounds is the box itself.
func (b Box) Bounds() Box { return b }
