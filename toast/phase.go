package toast

// Leaving is whether t has begun to fade out.
func (t *Toast) Leaving() bool { return t.Age >= t.Life }

// Gone is whether t has finished fading out and can be dropped.
func (t *Toast) Gone() bool { return t.Age >= t.Life+ExitFor }

// ease is a cubic ease-out of x in 0..1: quick at first, then settling.
func ease(x float64) float64 {
	x = min(1, max(0, x))
	return 1 - (1-x)*(1-x)*(1-x)
}

// Opacity is how opaque t is, 0 to 1: rising while it enters, 1 while it
// holds, falling while it leaves.
func (t *Toast) Opacity() float64 {
	if t.Leaving() {
		return 1 - min(1, (t.Age-t.Life)/ExitFor)
	}
	return ease(t.Age / EnterFor)
}

// Shift is how far t is from its place, in slides: -1 is one slide above it,
// as it starts to enter, 0 while it holds, and a little above again as it
// leaves.
func (t *Toast) Shift() float64 {
	if t.Leaving() {
		return -0.4 * ease((t.Age-t.Life)/ExitFor)
	}
	return ease(t.Age/EnterFor) - 1
}
