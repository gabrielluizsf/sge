package farm

// AnimSleep is a person asleep in a bed, lying on their back and breathing;
// it loops, and is shown whichever way the person faces.
const AnimSleep = "sleep"

// SleepArt is what an Art may also have whose bed is drawn with the person
// asleep in it, one picture, so the bed hides the body where it should.
// While they sleep, the bed's own tile is taken away.
type SleepArt interface {
	SleepBaked() bool
}

// sleepBaked reports whether a draws the bed with the sleeper in it.
func sleepBaked(a Art) bool {
	s, ok := a.(SleepArt)
	return ok && s.SleepBaked()
}

// Sleeper is what a Body may also have: whether its frames of the person
// asleep are made yet. A body without it always has them.
type Sleeper interface {
	Sleeps() bool
}

// sleeps reports whether the person b draws has frames asleep.
func sleeps(b Body) bool {
	s, ok := b.(Sleeper)
	return !ok || s.Sleeps()
}

// AnimSleepLeft and AnimSleepRight are a person asleep on a side of a
// double bed, the left and the right of it as the game draws it. A look
// that has no pictures of them draws them as AnimSleep.
const (
	AnimSleepLeft  = "sleep_l"
	AnimSleepRight = "sleep_r"
)

// SideSleeper is what a Body may also have, when it sleeps on one side of a
// double bed in pictures of its own: whether it has those of anim.
type SideSleeper interface {
	SleepsAs(anim string) bool
}

// sleepsAs reports whether b lies asleep as anim says.
func sleepsAs(b Body, anim string) bool {
	if s, ok := b.(SideSleeper); ok && anim != AnimSleep {
		return s.SleepsAs(anim)
	}
	return sleeps(b)
}
