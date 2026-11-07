package farm

import "github.com/gabrielluizsf/sge/pathfind"

// next is the villager's next job: to bed at night, their trade's work by
// day, a stroll about town otherwise.
func (m *mind) next(f *farm) *job {
	if m.v.guest != nil {
		return m.visiting(f)
	}
	switch m.v.now(f) {
	case doSleep:
		return m.home(0)
	case doWork:
		if j := m.work(f); j != nil && f.inSight(f.center(j.to)) {
			return j // work out of sight would only walk them off the screen
		}
	}
	return m.stroll(f)
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

// stroll is a walk to somewhere in sight along the street or on the
// plaza, and a rest there, or nil when there is none; a wanderer goes home now and then.
func (m *mind) stroll(f *farm) *job {
	v := m.v
	if m.walks++; v.role == roleWanderer && m.walks%5 == 0 {
		if j := m.home(10 + 15*v.luck.Float()); j != nil {
			return j
		}
	}
	c := v.place().Center
	for range strollTries {
		to := pathfind.Cell{Col: c.Col + v.luck.Between(-9, 9), Row: c.Row}
		if v.luck.Float() < 0.4 {
			to = pathfind.Cell{Col: c.Col + v.luck.Between(-2, 2), Row: c.Row + v.luck.Between(-1, 1)}
		}
		if f.inSight(f.center(to)) {
			return &job{to: to, slot: -1, rest: strollRest * (0.3 + v.luck.Float())}
		}
	}
	return nil
}

// strollTries is how many places a villager considers for a stroll before
// giving up, and strollRest the seconds they stand about after one, at
// most.
const (
	strollTries = 4
	strollRest  = 1.2
)
