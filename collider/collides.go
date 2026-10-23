package collider

// Collides reports whether a and b overlap, picking the exact test for the
// pair of shapes. A shape it does not know is tested by its bounds.
func Collides(a, b Collider) bool {
	switch a := a.(type) {
	case Circle:
		switch b := b.(type) {
		case Circle:
			return CircleCircle(a, b)
		case Box:
			return CircleBox(a, b)
		default:
			return CircleBox(a, b.Bounds())
		}
	case Box:
		if c, ok := b.(Circle); ok {
			return CircleBox(c, a)
		}
	}
	if c, ok := b.(Circle); ok {
		return CircleBox(c, a.Bounds())
	}
	return BoxBox(a.Bounds(), b.Bounds())
}
