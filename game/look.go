package game

import "github.com/gabrielluizsf/sge/sprite"

// Look is the sprite the object shows now: its animation's frame, or its
// sprite.
func (o *Object) Look() *sprite.Sprite {
	if o.Animator != nil {
		if f := o.Animator.Frame(); f != nil {
			return f
		}
	}
	return o.Sprite
}
