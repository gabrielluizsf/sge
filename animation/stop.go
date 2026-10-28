package animation

// Stop freezes the animation on the frame it is showing.
func (a *Animator) Stop() { a.playing = false }
