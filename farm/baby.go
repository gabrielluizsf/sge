package farm

import (
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/random"
)

// birth is a child being born: their parents' ids in the family tree,
// mother or father first, each parent's five numbers, and the day.
type birth struct {
	parents [2]string
	genes   [2][5]int
	day     int
}

// baby is the child of b, born: a girl or a boy as fate falls, named, each
// of their numbers one of their parents', in the family tree, a newborn
// with a span drawn by fate. Where they live, and with whom, is for those
// who bear them to say.
func (f *farm) baby(b birth) *villager {
	woman := f.fate.Float() < 0.5
	names := menNames
	if woman {
		names = womenNames
	}
	var genes [5]int
	for i := range genes {
		genes[i] = b.genes[f.fate.Intn(2)][i]
	}
	v := &villager{
		id: f.newID(), name: names[f.fate.Intn(len(names))], woman: woman, role: roleWanderer,
		who:  Who{Woman: woman, Role: roleWanderer, Genes: genes},
		kit:  &kit{money: 0, axe: 0, slot: slotCan},
		luck: random.New(uint32(f.fate.Intn(1 << 30))), wait: orderPause,
		life: calendar.Newborn(f.fate, b.day), inside: true,
	}
	v.mind = &mind{v: v}
	f.family.Add(lineage.Person{ID: v.id, Parents: b.parents, Born: b.day})
	f.age(v, b.day)
	return v
}
