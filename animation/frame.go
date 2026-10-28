package animation

import "github.com/gabrielluizsf/sge/sprite"

// Frame is the sprite the current animation is showing, or nil when
// there is none.
func (a *Animator) Frame() *sprite.Sprite {
	anim := a.anims[a.current]
	if anim == nil || len(anim.Frames) == 0 {
		return nil
	}
	return anim.Frames[a.FrameIndex()]
}
