package kit

import (
	"slices"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// Grown is a villager's look at a stage of life, as their sprites are
// baked: their age then, and the motions they have at it.
type Grown struct {
	Stage   life.Stage
	Years   float64
	Motions []human.Motion
}

// young are the motions of one too young to work.
var young = []human.Motion{human.Stand, human.Walk, human.Swim, human.Tread, human.Dive, human.Rise}

// Stages are each villager's stages of life, baked: a baby, standing; a
// child and a teen, playing; a grown-up, at their own age, and an elder,
// working.
func (v Villager) Stages() []Grown {
	own := max(v.Years(), life.AdultFrom+5)
	if v.Years() >= life.ElderFrom {
		own = 40
	}
	return []Grown{
		{life.Baby, 1, []human.Motion{human.Stand}},
		{life.Child, 8, v.childMotions()},
		{life.Teen, 15, v.childMotions()},
		{life.Adult, min(own, life.ElderFrom-1), v.Motions()},
		{life.Elder, 74, v.Motions()},
	}
}

// Years are how old the villager is as the game starts.
func (v Villager) Years() float64 { return v.Body.Years() }

// At is the villager's body at years: the same person, older or younger.
func (v Villager) At(years float64) human.Body {
	b := v.Body
	b.Age = age(years)
	return b
}

// childMotions are a child's and a teen's motions: the young's, and sleeping in their
// own bed, the one a household gives its child, unless they are the
// player, who has a bed of their own.
func (v Villager) childMotions() []human.Motion {
	if v.ID == "player" {
		return young
	}
	return append(slices.Clone(young), human.Sleeping...)
}
