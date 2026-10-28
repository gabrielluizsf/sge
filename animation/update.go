package animation

// Update moves the playing animation on by dt seconds.
func (a *Animator) Update(dt float64) {
	if a.playing && dt > 0 {
		a.elapsed += dt
	}
}
