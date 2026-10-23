package collider

// Offset is c moved by dx, dy: a collider given relative to an object,
// placed where the object is. A shape it does not know is moved by its
// bounds.
func Offset(c Collider, dx, dy float64) Collider {
	switch c := c.(type) {
	case Box:
		return Box{c.X + dx, c.Y + dy, c.Width, c.Height}
	case Circle:
		return Circle{c.X + dx, c.Y + dy, c.Radius}
	case nil:
		return nil
	}
	b := c.Bounds()
	return Box{b.X + dx, b.Y + dy, b.Width, b.Height}
}
