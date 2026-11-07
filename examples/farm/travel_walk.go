package farm

import "github.com/gabrielluizsf/sge/game"

// roadSnap is how many cells off the road's point a traveller may be
// before they are put on it, as on a first frame or after a long wait.
const roadSnap = 6

// travel is a villager on the road to or from the next town: they walk to
// the point of the road the plan has them at.
func (m *mind) travel(f *farm, o *game.Object, w whereabouts) intent {
	v := m.v
	v.inside, o.Hidden, v.homeFor = false, false, 0
	m.job, m.route, m.rest = nil, nil, 0
	target := f.point(w.road, w.at)
	d := target.Sub(feet(o))
	switch n := d.Length(); {
	case n > roadSnap*tile:
		o.Transform.Position = target
		return intent{}
	case n < 3:
		return intent{}
	}
	return intent{dir: d, upTo: d.Length()}
}
