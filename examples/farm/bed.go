package farm

import "github.com/gabrielluizsf/sge/lineage"

// goToBed is the player going to bed: married, late enough, and able to,
// they are only asked whether to have a child, and the night passes as
// they answer; otherwise they sleep, if it is late enough.
func (f *farm) goToBed() {
	if f.clock.Between(bedtime, wakeAt) && f.mayHaveChild() {
		f.question = childQuestion()
		f.sound("ui_open", nil)
		return
	}
	f.sleep()
}

// mayHaveChild reports whether the player may decide on a child: married
// to a spouse alive, not expecting one already, the woman of the two
// young enough.
func (f *farm) mayHaveChild() bool {
	s := f.self.spouse
	if s == nil || !s.bond.MayHaveChild() || f.self.expecting != nil {
		return false
	}
	mother := s.life
	if f.self.woman {
		mother = f.self.life
	}
	return yearsOf(mother, f.clock.Day) < fertileUntil
}

// conceive is the player and their spouse deciding on a child tonight,
// born lineage.Gestation days on.
func (f *farm) conceive() {
	if !f.mayHaveChild() {
		return
	}
	e := lineage.Expect([2]string{f.self.id, f.self.spouse.id}, f.clock.Day)
	f.self.expecting = &e
}
