package collider

// CircleBox reports whether a circle overlaps a box, by measuring from the
// circle's centre to the nearest point of the box.
func CircleBox(c Circle, b Box) bool {
	dx := c.X - clamp(c.X, b.X, b.X+b.Width)
	dy := c.Y - clamp(c.Y, b.Y, b.Y+b.Height)
	return dx*dx+dy*dy <= c.Radius*c.Radius
}
