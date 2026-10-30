package farm

import "github.com/gabrielluizsf/sge/pathfind"

// next is the villager's next job: to bed at night, their trade's work by
// day, a stroll about town otherwise.
func (m *mind) next(f *farm) *job {
	switch m.v.now(f) {
	case doSleep:
		return m.home(0)
	case doWork:
		if j := m.work(f); j != nil {
			return j
		}
	}
	return m.stroll()
}

// work is the job of the villager's trade, or nil when there is none to
// do now.
func (m *mind) work(f *farm) *job {
	switch m.v.role {
	case roleFarmer:
		return m.farmWork(f)
	case roleLumberjack:
		return m.lumberWork(f)
	case roleSmith:
		return m.smithWork(f)
	}
	return nil
}

// home is going home and staying in for stay seconds, or all night.
func (m *mind) home(stay float64) *job {
	v := m.v
	h := v.home
	if h == nil {
		return nil
	}
	face := "down"
	if h.Row < h.Door.Row {
		face = "up"
	}
	return &job{to: pathfind.Cell(h.Door), face: face, slot: -1, then: func(*farm) { v.inside, v.homeFor = true, stay }}
}

// stroll is a walk to somewhere along the street or on the plaza, and a
// rest there; a wanderer goes home now and then.
func (m *mind) stroll() *job {
	v := m.v
	if m.walks++; v.role == roleWanderer && m.walks%5 == 0 {
		if j := m.home(10 + 15*v.luck.Float()); j != nil {
			return j
		}
	}
	c := v.town.Center
	to := pathfind.Cell{Col: c.Col + v.luck.Between(-9, 9), Row: c.Row}
	if v.luck.Float() < 0.4 {
		to = pathfind.Cell{Col: c.Col + v.luck.Between(-2, 2), Row: c.Row + v.luck.Between(-1, 1)}
	}
	return &job{to: to, slot: -1, rest: 2 + 4*v.luck.Float()}
}
