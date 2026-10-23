package animation

// Done reports whether an animation that does not loop has reached its
// last frame's end. A looping one is never done.
func (a *Animator) Done() bool {
	anim := a.anims[a.current]
	if anim == nil || anim.Loop || anim.FPS <= 0 {
		return false
	}
	return a.elapsed*anim.FPS >= float64(len(anim.Frames))
}
