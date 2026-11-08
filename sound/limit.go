package sound

// cutFade is how many seconds a sound takes to fade out when another of its
// name starts: too quick to hear as a fade, too slow to click.
const cutFade = 0.03

// limitKnee is the loudness up to which a sample passes unchanged.
const limitKnee = 0.7

// limit is sample s kept from -1 to 1: unchanged while it is quiet, and above
// the knee bending smoothly toward the limit, never reaching past it.
func limit(s float32) float32 {
	sign := float32(1)
	if s < 0 {
		s, sign = -s, -1
	}
	if s <= limitKnee {
		return sign * s
	}
	over := (s - limitKnee) / (1 - limitKnee)
	return sign * (limitKnee + (1-limitKnee)*over/(1+over))
}
