package animation

// Play starts the named animation from its first frame, unless it is the
// one already playing, which carries on. It reports whether there is an
// animation of that name.
func (a *Animator) Play(name string) bool {
	if _, ok := a.anims[name]; !ok {
		return false
	}
	if a.current == name && a.playing {
		return true
	}
	a.current, a.elapsed, a.playing = name, 0, true
	return true
}
