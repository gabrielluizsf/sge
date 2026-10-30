package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// mind is a villager's controller: it picks a job as the hour and their
// trade say, walks the way pathfind finds to where it is done, does it,
// and rests a while before the next.
type mind struct {
	v     *villager
	job   *job
	route *pathfind.Route
	rest  float64
	walks int
	// stuck is how long they have been walking without getting anywhere,
	// since last.
	stuck float64
	last  vector.Vector3
}

// job is a place to stand, the way to face there, what to hold (-1 keeps
// what is held), how many times to use it, what then happens, and how
// long to rest after.
type job struct {
	to   pathfind.Cell
	face string
	slot int
	uses int
	then func(f *farm)
	rest float64
}

// giveUp is how long a villager walks without getting anywhere before
// giving the job up.
const giveUp = 3

func (m *mind) control(f *farm, o *game.Object, dt float64) intent {
	v := m.v
	if v.inside {
		if v.homeFor -= dt; v.homeFor > 0 || v.now(f) == doSleep {
			o.Hidden = true
			return intent{}
		}
		v.inside, o.Hidden = false, false
	}
	if m.rest > 0 {
		m.rest -= dt
		return intent{}
	}
	if m.job == nil {
		if m.job = m.next(f); m.job == nil {
			m.rest = 1
			return intent{}
		}
	}
	j := m.job
	at := feet(o)
	if m.route == nil && !m.plan(f, at) {
		return intent{}
	}
	if way, done := m.route.Steer(at, f.center); !done {
		if m.stuck += dt; at.Sub(m.last).Length() > 4 {
			m.stuck, m.last = 0, at
		}
		if m.stuck > giveUp {
			m.drop(1)
			return intent{}
		}
		return intent{dir: way, upTo: way.Length()}
	}
	in := intent{face: j.face}
	if j.slot >= 0 {
		v.kit.slot = j.slot
	}
	if j.uses > 0 {
		j.uses--
		in.use = true
		return in
	}
	if j.then != nil {
		j.then(f)
	}
	m.drop(j.rest)
	return in
}

// plan finds the way from at to the job's place, and reports whether
// there is one; a job with none is given up.
func (m *mind) plan(f *farm, at vector.Vector3) bool {
	col, row := f.cellOf(at)
	way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, m.job.to)
	if way == nil {
		m.drop(1 + m.v.luck.Float()*2)
		return false
	}
	if len(way) == 0 {
		way = []pathfind.Cell{m.job.to} // to the middle of the cell they are on
	}
	m.route = &pathfind.Route{Cells: way}
	m.stuck, m.last = 0, at
	return true
}

// drop ends the job, resting rest seconds.
func (m *mind) drop(rest float64) {
	m.job, m.route, m.rest = nil, nil, rest
}

// grid is where people walk: every cell they need not swim across.
func (f *farm) grid() pathfind.Grid {
	return pathfind.Grid{Open: f.open, Limit: 2500}
}

// now is what v's schedule says to do at this hour.
func (v *villager) now(f *farm) string { return schedules[v.role].At(f.clock) }
