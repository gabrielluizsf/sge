package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/life"
)

// age brings villager v's years up to day: their whole years, and their
// time of life, their body drawn anew when it changes, if it has been
// drawn already.
func (f *farm) age(v *villager, day int) {
	y := yearsOf(v.life, day)
	st := life.StageOf(y)
	v.who.Years = int(y)
	if st == v.who.Stage {
		return
	}
	v.who.Stage = st
	if v.body == nil {
		return
	}
	v.body = f.art.Person(v.who)
	if v.obj != nil {
		dress(v.obj, v.body)
	}
	for _, r := range f.residents {
		if r.v == v {
			dress(r.obj, v.body)
		}
	}
}

// growOld ages every villager to day, and those whose span has run out die
// of old age.
func (f *farm) growOld(day int) {
	for _, v := range append([]*villager(nil), f.people...) {
		if calendar.Gone(v.life, day) {
			f.bury(v)
			continue
		}
		f.age(v, day)
	}
}

// dress gives person o the clips of body b, in place of those they had.
func dress(o *game.Object, b Body) {
	for _, dir := range Facings(b.Facings()) {
		addMoves(o.Animator, b, dir)
	}
}
