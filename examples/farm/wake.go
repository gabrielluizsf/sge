package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// wake is a town's people's objects, as the keeper wakes it: settled the
// first time, each standing where they were when it last slept, moved by
// their mind and kept there as they walk; not those who live with the
// player now.
func (f *farm) wake(t *town.Town) []*game.Object {
	f.settle(t)
	f.catchUp(t)
	var objs []*game.Object
	var later []*sprite.Sprite // the frames to make once all have their first
	for _, v := range f.people {
		if !f.isIn(v, t) {
			continue
		}
		o := f.spawn(v, f.restAt(v, t))
		objs = append(objs, o)
		later = append(later, bodyFrames(v.body)...)
	}
	f.warm.Add(later...)
	return objs
}

// spawn is villager v's object, standing at pos and driven by their mind.
func (f *farm) spawn(v *villager, pos vector.Vector3) *game.Object {
	o, p := f.personOf(v)
	p.ctl, v.hands = v.mind, p
	o.Transform.Position, v.pos = pos, pos
	o.Hidden = v.inside
	o.AddScript(game.ScriptFunc(func(o *game.Object, _ float64) { v.pos = o.Transform.Position }))
	v.obj = o
	f.warm.Add(bodyFirst(v.body))
	return o
}

// personOf is an object of villager v standing still, as they are seen
// at home; wake gives the one outside their mind.
func (f *farm) personOf(v *villager) (*game.Object, *farmer) {
	if v.body == nil {
		f.age(v, f.clock.Day)
		v.body = f.art.Person(v.who)
	}
	return newPerson(f, person{name: "villager", body: v.body, kit: v.kit, ctl: still{}})
}

// still is a controller that keeps one where they stand.
type still struct{}

func (still) control(*farm, *game.Object, float64) intent { return intent{} }
