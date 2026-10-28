package collider

// BoxBox reports whether two boxes overlap. Boxes that only touch along an
// edge do not.
func BoxBox(a, b Box) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}
