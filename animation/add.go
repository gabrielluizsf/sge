package animation

// Add names an animation, replacing any of the same name.
func (a *Animator) Add(name string, anim *Animation) {
	a.anims[name] = anim
}
