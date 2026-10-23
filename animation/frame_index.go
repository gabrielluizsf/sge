package animation

// FrameIndex is which frame of the current animation is showing.
func (a *Animator) FrameIndex() int {
	anim := a.anims[a.current]
	if anim == nil || len(anim.Frames) == 0 || anim.FPS <= 0 {
		return 0
	}
	i := int(a.elapsed * anim.FPS)
	if anim.Loop {
		return i % len(anim.Frames)
	}
	return min(i, len(anim.Frames)-1)
}
