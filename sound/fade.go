package sound

// Fade eases the loop's gain to gain, over a fifth of a second or so.
func (l *Loop) Fade(gain float64) {
	l.m.mu.Lock()
	defer l.m.mu.Unlock()
	l.target = gain
}

// fadeTime is about how long a loop takes to ease to a gain, in seconds.
const fadeTime = 0.2
