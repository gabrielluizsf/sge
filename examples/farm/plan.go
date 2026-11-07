package farm

// chore is one thing a villager spends part of a day on.
type chore int

// The chores of a day: what the player does, and the trip to the next
// town to trade there.
const (
	choreChop chore = iota
	choreMine
	chorePlant
	choreForge
	choreTrade
	choreVisit
	choreCount
)

// spell is a stretch of a day, from one hour to another, spent on a chore.
type spell struct {
	from, to float64
	doing    chore
}

// dayPlan is a villager's day: their spells in the order of the hours.
type dayPlan []spell

// at is the spell that hour falls in, if any.
func (p dayPlan) at(hour float64) (spell, bool) {
	for _, s := range p {
		if hour >= s.from && hour < s.to {
			return s, true
		}
	}
	return spell{}, false
}

// visit is the spell of the day spent in the next town, if there is one.
func (p dayPlan) visit() (spell, bool) {
	for _, s := range p {
		if s.doing == choreVisit {
			return s, true
		}
	}
	return spell{}, false
}

// index is the position of the spell that has the hour, or -1 when none has.
func (p dayPlan) index(hour float64) int {
	for i, s := range p {
		if hour >= s.from && hour < s.to {
			return i
		}
	}
	return -1
}
