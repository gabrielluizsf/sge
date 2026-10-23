package game

import "github.com/gabrielluizsf/sge/collider"

// Collides reports whether the object touches other. Objects without a
// collider touch nothing.
func (o *Object) Collides(other *Object) bool {
	if o.Collider == nil || other.Collider == nil || o == other {
		return false
	}
	return collider.Collides(o.WorldCollider(), other.WorldCollider())
}
