package game

import "github.com/gabrielluizsf/sge/collider"

// WorldCollider is the object's collider where the object is.
func (o *Object) WorldCollider() collider.Collider {
	p := o.Transform.Position
	return collider.Offset(o.Collider, p.X, p.Y)
}
