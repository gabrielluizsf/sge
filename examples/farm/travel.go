package farm

// stay is where a villager is in the course of a visit to the next town.
type stay int

// The course of a visit: at home, walking out, there, walking back.
const (
	atHome stay = iota
	outbound
	there
	inbound
)

// whereabouts is where a villager is by the day's plan: the stay, how far
// along the road to the next town, from 0 to 1, and the road.
type whereabouts struct {
	stay stay
	at   float64
	road *road
}

// walking reports whether the villager is on the road.
func (w whereabouts) walking() bool { return w.stay == outbound || w.stay == inbound }

// whereabouts is where v is now, by their plan alone: nothing is
// simulated for it, so it costs the same however far away they are.
func (f *farm) whereabouts(v *villager) whereabouts {
	return f.whereaboutsAt(v, f.clock.Day, f.clock.Hour)
}

// whereaboutsAt is where v is at an hour of any day, by that day's plan
// alone.
func (f *farm) whereaboutsAt(v *villager, day int, hour float64) whereabouts {
	plan := f.planOf(v)
	if day != f.clock.Day {
		plan = f.planOn(v, day)
	}
	s, ok := plan.visit()
	if !ok {
		return whereabouts{}
	}
	r := f.roadFrom(v.town)
	if r == nil {
		return whereabouts{}
	}
	h, trip := hour, r.hours
	switch {
	case h < s.from || h >= s.to:
		return whereabouts{road: r}
	case h < s.from+trip:
		return whereabouts{outbound, (h - s.from) / trip, r}
	case h < s.to-trip:
		return whereabouts{there, 1, r}
	}
	return whereabouts{inbound, (s.to - h) / trip, r}
}
