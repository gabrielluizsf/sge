package farm

import "github.com/gabrielluizsf/sge/town"

// The most hours of a town's time away that are settled when it wakes, and
// the least: a longer time is settled as the last of these hours, so a
// town's wake costs at most that much work however long it was left; a
// shorter one is not worth settling.
const (
	maxAway = 3 * 24.0
	minAway = 0.5
)

// hours is the time of the game in hours since its first day began.
func (f *farm) hours() float64 { return float64(f.clock.Day)*24 + f.clock.Hour }

// leftAt takes note of when town t was last seen, as it goes to sleep.
func (f *farm) leftAt(t *town.Town) {
	if f.away == nil {
		f.away = map[*town.Town]float64{}
	}
	f.away[t] = f.hours()
}

// catchUp settles what the people of town t did while no one saw them, from
// the plans they kept: the days since it was last seen, the last maxAway
// hours at most, one after another. A town seen for the first time has
// nothing to settle.
func (f *farm) catchUp(t *town.Town) {
	since, ok := f.away[t]
	if !ok {
		return
	}
	delete(f.away, t)
	now := f.hours()
	if now-since < minAway {
		return
	}
	f.settleSpan(t, max(since, now-maxAway), now)
}
