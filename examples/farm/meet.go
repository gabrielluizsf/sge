package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
)

// meetFor is how long, in seconds, two who meet stand facing each other.
const meetFor = 2.0

// canMeet reports whether v may go to o to trade: o is a grown villager
// in the same town, in sight, out and about, with something to trade with v.
func (f *farm) canMeet(v, o *villager) bool {
	if o == v || o.ours || o.obj == nil || o.obj.Hidden || o.inside || o.stage() < life.Adult || o.place() != v.place() {
		return false
	}
	if o.mind.hold > 0 || f.whereabouts(o).walking() || !f.inSight(feet(o.obj)) {
		return false
	}
	_, ok := saleBetween(v, o)
	return ok
}

// meetJob is a walk to another villager in sight with whom there is
// something to trade, to stand beside them: they meet there.
func (m *mind) meetJob(f *farm) *job {
	v := m.v
	var near []*villager
	for _, o := range f.people {
		if f.canMeet(v, o) {
			near = append(near, o)
		}
	}
	if len(near) == 0 {
		return nil
	}
	o := near[v.luck.Intn(len(near))]
	col, row := f.cellOf(feet(o.obj))
	for _, side := range [2]struct {
		dc   int
		face string
	}{{-1, "right"}, {1, "left"}} {
		to := pathfind.Cell{Col: col + side.dc, Row: row}
		if f.open(to.Col, to.Row) && f.inSight(f.center(to)) {
			return &job{to: to, face: side.face, slot: -1, rest: 1, then: func(f *farm) { f.meet(v, o) }}
		}
	}
	return nil
}

// meet is villagers v and o, who have come together, stopping to face each
// other and trade: a ware for its price, with a note over each.
func (f *farm) meet(v, o *villager) {
	if v.obj == nil || o.obj == nil || feet(v.obj).Sub(feet(o.obj)).Length() > 3*tile {
		return
	}
	v.mind.hold, v.mind.facing = meetFor, faceTo(feet(v.obj), o.obj)
	o.mind.hold, o.mind.facing = meetFor, faceTo(feet(o.obj), v.obj)
	o.mind.job, o.mind.route = nil, nil
	s, ok := saleBetween(v, o)
	if !ok {
		return
	}
	s.settle()
	f.float(feet(s.from.obj), i18n.T("+$%d", s.ware.price), 'y')
	f.float(feet(s.to.obj), i18n.T(s.ware.label), 'r')
	at := feet(s.to.obj)
	f.sound("coin", &at)
}
