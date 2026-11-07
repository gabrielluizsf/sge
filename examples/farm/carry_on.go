package farm

import (
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// carryOn is where v of town t stands as it wakes when they are at a chore
// by the plan at this hour: at the place of it, with its job in hand, so
// that those watching see them go on with it, not begin it from their door.
// It reports false when they have no chore now, or are not in their town.
func (f *farm) carryOn(v *villager, t *town.Town) (vector.Vector3, bool) {
	if v.town != t || v.ours || v.stage() < life.Adult || v.now(f) != doWork {
		return vector.Vector3{}, false
	}
	p := f.presenceAt(v, f.clock.Day, f.clock.Hour)
	if p.job == nil {
		return vector.Vector3{}, false
	}
	v.inside = false
	v.mind.job, v.mind.route, v.mind.rest = p.job, nil, 0
	return f.center(p.job.to), true
}
