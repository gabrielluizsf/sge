package collider

// Collides reports whether a and b overlap, picking the exact test for the
// pair of shapes. The shapes it knows are tested first, as most are those;
// an Overlapper decides for itself, and another shape is tested by its
// bounds.
func Collides(a, b Collider) bool {
	switch a := a.(type) {
	case Circle:
		switch b := b.(type) {
		case Circle:
			return CircleCircle(a, b)
		case Box:
			return CircleBox(a, b)
		}
	case Box:
		switch b := b.(type) {
		case Circle:
			return CircleBox(b, a)
		case Box:
			return BoxBox(a, b)
		}
	}
	if o, ok := a.(Overlapper); ok {
		return o.Overlaps(b)
	}
	if o, ok := b.(Overlapper); ok {
		return o.Overlaps(a)
	}
	if c, ok := a.(Circle); ok {
		return CircleBox(c, b.Bounds())
	}
	if c, ok := b.(Circle); ok {
		return CircleBox(c, a.Bounds())
	}
	return BoxBox(a.Bounds(), b.Bounds())
}
