package game

import "github.com/gabrielluizsf/sge/collider"

// WorldCollider is the object's collider where the object is. It is
// worked out again only when the position or the collider changed.
func (o *Object) WorldCollider() collider.Collider {
	p := o.Transform.Position
	if o.world == nil || p.X != o.worldX || p.Y != o.worldY || !sameShape(o.Collider, o.worldFrom) {
		o.world = collider.Offset(o.Collider, p.X, p.Y)
		o.worldFrom, o.worldX, o.worldY = o.Collider, p.X, p.Y
	}
	return o.world
}
