package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// trafficEvery is how often, in seconds, the villagers on visits are
// looked at.
const trafficEvery = 0.5

// farFromHome is how many cells from their own town's middle someone
// saved elsewhere is put back at their door, when the plan has them home.
const farFromHome = 40

// isIn reports whether v's object belongs in town t as it wakes: they are
// at home in it by the plan, or they visit it now.
func (f *farm) isIn(v *villager, t *town.Town) bool {
	if v.ours {
		return false
	}
	w := f.whereabouts(v)
	switch {
	case v.town == t:
		return w.stay == atHome
	case w.stay == there:
		return w.road.to == t
	}
	return false
}

// restAt is where v stands as t wakes: where they were left, but a
// visitor on the plaza they arrived at, and one who is home but was left
// elsewhere at their door.
func (f *farm) restAt(v *villager, t *town.Town) vector.Vector3 {
	w := f.whereabouts(v)
	if w.stay == there {
		v.guest = t
		v.inside = false
		return f.center(pathCell(t.Center))
	}
	v.guest = nil
	if at, ok := f.carryOn(v, t); ok {
		return at
	}
	if v.home != nil && v.town == t {
		d := f.center(pathCell(t.Center)).Sub(v.pos)
		if d.Length() > farFromHome*tile {
			return f.center(pathCell(v.home.Door))
		}
	}
	return v.pos
}

// traffic puts the villagers on visits where their plans have them, once
// in a while: those on the road in sight walk it, the guests of an awake
// town are in it, and those who left, or came back, to an awake town are
// not or are there. A town that is asleep is not looked at, and nothing is
// done for those out of range.
func (f *farm) traffic(dt float64) {
	if f.trafficIn -= dt; f.trafficIn > 0 || f.room != nil || f.keeper == nil || f.player.Game() == nil {
		return
	}
	f.trafficIn = trafficEvery
	for _, v := range f.people {
		if v.ours || v.town == nil || !f.keeper.Awake(v.town) && !f.awakeRoad(v.town) {
			continue
		}
		if v.stage() < life.Adult {
			continue
		}
		f.keep(v, f.whereabouts(v))
	}
}

// awakeRoad reports whether the town t's road leads to a town that is awake.
func (f *farm) awakeRoad(t *town.Town) bool {
	r := f.roadFrom(t)
	return r != nil && f.keeper.Awake(r.to)
}

// keep makes villager v's object as w says it should be: there when they
// are in an awake town, or on the road in sight, and not otherwise.
func (f *farm) keep(v *villager, w whereabouts) {
	if w.road == nil {
		return
	}
	var at vector.Vector3
	here := false
	switch w.stay {
	case atHome:
		v.guest = nil
		if s, ok := f.planOf(v).visit(); !ok || f.clock.Hour < s.from {
			return // not gone yet: the town's own wake and mind have them
		}
		at, here = f.center(pathCell(v.town.Center)), f.keeper.Awake(v.town)
	case there:
		v.guest = w.road.to
		v.inside = false
		at, here = f.center(pathCell(w.road.to.Center)), f.keeper.Awake(w.road.to)
	default:
		v.guest = nil
		at = f.point(w.road, w.at)
		here = (f.keeper.Awake(v.town) || f.keeper.Awake(w.road.to)) && f.inSight(at)
	}
	switch {
	case v.obj != nil && !here:
		v.obj.Destroy()
		v.obj = nil
	case v.obj == nil && here:
		f.player.Game().Add(f.spawn(v, at))
	}
}

// leaves is a villager's object let go as their town sleeps, whether the
// keeper made it or traffic.
func leaves(o *game.Object) {
	if o != nil {
		o.Destroy()
	}
}
