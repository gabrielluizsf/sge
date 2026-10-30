package farm

// days lives through each day the clock has reached since the farm last
// looked, in turn.
func (f *farm) days() {
	for f.seen < f.clock.Day {
		f.seen++
		f.newDay(f.seen)
	}
}

// newDay is the morning of day: the engaged wed, children due are born,
// couples in the towns may decide on one, everyone is a day older, and
// those whose time is up die, the player too, an heir taking over.
func (f *farm) newDay(day int) {
	f.weddings(day)
	f.births(day)
	f.townBirths(day)
	f.growOld(day)
	f.playerAges(day)
}

// passDays passes n days at once, as if slept through, each lived in
// turn: what F9 does a year at a time, and the tests too.
func (f *farm) passDays(n int) {
	f.clock.Day += n
	f.days()
}
