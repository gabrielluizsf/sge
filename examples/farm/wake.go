package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// wake is a town's people's objects, as the keeper wakes it: settled the
// first time, each standing where they were when it last slept, moved by
// their mind and kept there as they walk; not those who live with the
// player now.
func (f *farm) wake(t *town.Town) []*game.Object {
	f.settle(t)
	var objs []*game.Object
	for _, v := range f.people {
		if v.town != t || v.ours {
			continue
		}
		o, p := f.personOf(v)
		p.ctl = v.mind
		o.Transform.Position = v.pos
		o.Hidden = v.inside
		o.AddScript(game.ScriptFunc(func(o *game.Object, _ float64) { v.pos = o.Transform.Position }))
		v.obj = o
		objs = append(objs, o)
	}
	return objs
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
