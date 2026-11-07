package farm

import "github.com/gabrielluizsf/sge/life"

// onDuty reports whether v is one of the fieldHands villagers of their
// trade who come to the farm today: the grown ones of that trade take
// turns, a pair a day.
func (f *farm) onDuty(v *villager) bool {
	n, rank := 0, -1
	for _, p := range f.people {
		if p.role != v.role || p.ours || p.stage() < life.Adult || p.town != nil && !f.keeper.Awake(p.town) {
			continue
		}
		if p == v {
			rank = n
		}
		n++
	}
	if rank < 0 {
		return false
	}
	return (rank-fieldHands*f.clock.Day%n+n)%n < fieldHands
}
