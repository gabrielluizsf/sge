package collider

// CircleCircle reports whether two circles overlap.
func CircleCircle(a, b Circle) bool {
	dx, dy, r := a.X-b.X, a.Y-b.Y, a.Radius+b.Radius
	return dx*dx+dy*dy <= r*r
}
