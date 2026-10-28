package main

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// wake is a town's people's objects, as the keeper wakes it: settled the
// first time, each standing where they were when it last slept, moved by
// their mind and kept there as they walk.
func (f *farm) wake(t *town.Town) []*game.Object {
	f.settle(t)
	var objs []*game.Object
	for _, v := range f.people {
		if v.town != t {
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
	if v.frames == nil {
		v.frames = art.PersonFrames(v.look)
	}
	frames := v.frames
	return newPerson(f, person{name: "villager", frame: func(n string) string { return frames[n] }, kit: v.kit, ctl: still{}})
}

// still is a controller that keeps one where they stand.
type still struct{}

func (still) control(*farm, *game.Object, float64) intent { return intent{} }
