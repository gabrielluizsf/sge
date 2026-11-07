package farm

import "github.com/gabrielluizsf/sge/pathfind"

// visiting is what a villager does in the town they are visiting: meet
// the people of it to trade, and stroll about its plaza between.
func (m *mind) visiting(f *farm) *job {
	if j := m.meetJob(f); j != nil {
		return j
	}
	return m.stroll(f)
}

// tradeJob is a villager's trade for the day: a meeting with someone to
// trade with, in sight, or else their town's stall.
func (m *mind) tradeJob(f *farm) *job {
	if j := m.meetJob(f); j != nil {
		return j
	}
	k := m.v.kit
	return m.atStall(func(f *farm, at pathfind.Cell) {
		if worth(k) > 0 {
			f.sell(k, f.center(at))
		}
	})
}
