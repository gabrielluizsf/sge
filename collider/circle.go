package collider

// Circle is a circle centred on X, Y.
type Circle struct {
	X, Y, Radius float64
}

// Bounds is the square the circle is drawn in.
func (c Circle) Bounds() Box {
	return Box{c.X - c.Radius, c.Y - c.Radius, 2 * c.Radius, 2 * c.Radius}
}
