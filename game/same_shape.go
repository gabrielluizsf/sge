package game

import "github.com/gabrielluizsf/sge/collider"

// sameShape reports whether a and b are known to be the same collider.
// Only the shapes sge knows are compared: comparing a collider of another
// type could panic, if it holds a slice say, so it is never the same.
func sameShape(a, b collider.Collider) bool {
	switch a := a.(type) {
	case collider.Box:
		b, ok := b.(collider.Box)
		return ok && a == b
	case collider.Circle:
		b, ok := b.(collider.Circle)
		return ok && a == b
	}
	return false
}
