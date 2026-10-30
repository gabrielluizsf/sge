package farm

import (
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/town"
)

// Couples in the towns have a child now and then: about one a year, while
// the mother is young enough, and none once a town has crowd people per
// house, so the towns never fill up.
const (
	childChance = 1.0 / daysPerYear
	crowd       = 2
)

// townBirths has the towns' couples decide on children on day, and their
// children born when due, at home in their mother's house.
func (f *farm) townBirths(day int) {
	for _, v := range append([]*villager(nil), f.people...) {
		switch {
		case v.expecting != nil && v.expecting.Born(day):
			f.bearInTown(v, day)
		case v.expecting == nil && f.mayExpect(v, day) && f.fate.Float() < childChance:
			e := lineage.Expect([2]string{v.id, v.partner.id}, day)
			v.expecting = &e
		}
	}
}

// mayExpect reports whether v, a woman of a town, may decide on a child
// on day: she has a partner, both grown, she young enough, at home in a
// town not yet crowded.
func (f *farm) mayExpect(v *villager, day int) bool {
	if !v.woman || v.ours || v.partner == nil || v.home == nil || v.town == nil {
		return false
	}
	if !grownOn(v.life, day) || !grownOn(v.partner.life, day) || yearsOf(v.life, day) >= fertileUntil {
		return false
	}
	return f.townSize(v.town) < crowd*len(v.town.Houses)
}

// bearInTown is v's child born on day, in her house.
func (f *farm) bearInTown(v *villager, day int) {
	e := v.expecting
	v.expecting = nil
	other := v.who.Genes
	if v.partner != nil {
		other = v.partner.who.Genes
	}
	b := f.baby(birth{parents: e.Parents, genes: [2][5]int{v.who.Genes, other}, day: day})
	b.town, b.home = v.town, v.home
	b.pos = f.ground.CellCenter(v.home.Door.Col, v.home.Door.Row)
	f.people = append(f.people, b)
}

// townSize is how many live in town t.
func (f *farm) townSize(t *town.Town) int {
	n := 0
	for _, v := range f.people {
		if v.town == t && !v.ours {
			n++
		}
	}
	return n
}
