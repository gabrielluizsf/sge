package animation

// NewAnimator is an animator with no animations, playing nothing.
func NewAnimator() *Animator {
	return &Animator{anims: map[string]*Animation{}}
}
