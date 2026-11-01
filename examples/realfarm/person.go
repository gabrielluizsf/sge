package main

import (
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/life"
)

// Person is who as the baker has baked them: the villager they are drawn
// as, at their stage of life, or at the nearest stage baked while the
// baker is still at work.
func (a *realArt) Person(who farm.Who) farm.Body {
	id := villagerOf(who)
	return &realBody{art: a, id: id, stage: a.stageBaked(id, who.Stage)}
}

// stageOrder are the stages of life by how near they are to each stage.
var stageOrder = map[life.Stage][]life.Stage{
	life.Baby:  {life.Baby, life.Child, life.Teen, life.Adult, life.Elder},
	life.Child: {life.Child, life.Teen, life.Baby, life.Adult, life.Elder},
	life.Teen:  {life.Teen, life.Child, life.Adult, life.Elder, life.Baby},
	life.Adult: {life.Adult, life.Teen, life.Elder, life.Child, life.Baby},
	life.Elder: {life.Elder, life.Adult, life.Teen, life.Child, life.Baby},
}

// stageBaked is the stage nearest to want that villager id has standing
// baked, want itself once it is.
func (a *realArt) stageBaked(id string, want life.Stage) life.Stage {
	for _, s := range stageOrder[want] {
		if a.shelf.Has(bakes.PersonName(id, s, "stand") + "_d0_f0") {
			return s
		}
	}
	return want
}
