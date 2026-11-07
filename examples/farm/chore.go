package farm

import "github.com/gabrielluizsf/sge/life"

// work is the job of what the plan has the villager on now, or of their
// trade when the plan has nothing, or nil when there is none to do.
func (m *mind) work(f *farm) *job {
	if s, ok := f.planOf(m.v).at(f.clock.Hour); ok {
		if j := m.chore(f, s.doing); j != nil {
			return j
		}
	}
	return m.tradeWork(f)
}

// chore is the job of doing c: felling, mining, tending a plot, working
// the forge, or trading; nil when there is none to do in sight.
func (m *mind) chore(f *farm, c chore) *job {
	v := m.v
	switch c {
	case choreChop:
		return m.chopJob(f)
	case choreMine:
		if v.stage() >= life.Adult {
			return m.mineJob(f)
		}
	case chorePlant:
		return m.farmWork(f)
	case choreForge:
		if v.smith() {
			return m.smithWork(f)
		}
	case choreTrade:
		return m.tradeJob(f)
	}
	return nil
}

// tradeWork is the work of the villager's trade.
func (m *mind) tradeWork(f *farm) *job {
	switch m.v.role {
	case roleFarmer:
		return m.farmWork(f)
	case roleLumberjack:
		return m.lumberWork(f)
	case roleSmith:
		return m.smithWork(f)
	case roleWanderer:
		return m.oddJobs(f)
	}
	return nil
}
