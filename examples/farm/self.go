package farm

import (
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
)

// self is who the player is, as one of a line of farmers: their id in the
// family tree, their name, a woman or a man, their five numbers, all 0 for
// one who founds the line, and their years; their spouse, once wed, and
// the child they are expecting, nil for none.
type self struct {
	id        string
	name      string
	woman     bool
	genes     [5]int
	life      life.Life
	spouse    *villager
	expecting *lineage.Expecting
}

// founderName is the name of the farmer who starts a line.
const founderName = "Joaquim"

// who is the player as the art sees them on day.
func (s *self) who(day int) Who {
	y := yearsOf(s.life, day)
	return Who{Player: true, Woman: s.woman, Genes: s.genes, Stage: life.StageOf(y), Years: int(y)}
}

// newID is a new id for the family tree, never given before.
func (f *farm) newID() string {
	f.ids++
	return "p" + itoa(f.ids)
}

// found is a new farmer founding a line on the farm: founderYears old
// today, in the family tree with no parents.
func (f *farm) found() *self {
	day := f.clock.Day
	s := &self{id: f.newID(), name: founderName, life: calendar.Draw(f.fate, founderYears, day)}
	f.family.Add(lineage.Person{ID: s.id, Born: s.life.Born})
	return s
}
