//go:build gots_overlay

// Overlay for the overlap checks, which Go makes on addresses as numbers:
// gots has no such addresses, but two slices of one array end it at the
// same element, and where each starts is told by how far that is.
//
// A slice whose capacity was cut by a full slice expression no longer ends
// its array there, so its overlaps go unseen: the checks catch misuse, and
// such a slice is simply taken as apart.

package alias

func AnyOverlap(x, y []byte) bool {
	if len(x) == 0 || len(y) == 0 {
		return false
	}
	xa, ya := x[:cap(x)], y[:cap(y)]
	if &xa[len(xa)-1] != &ya[len(ya)-1] {
		return false
	}
	// Each slice starts cap before the array's end.
	xs, ys := -cap(x), -cap(y)
	return xs < ys+len(y) && ys < xs+len(x)
}
