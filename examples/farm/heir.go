package farm

import "github.com/gabrielluizsf/sge/bond"

// playerAges ages the player to day: drawn anew as their time of life
// changes, and, their span run out, dead of old age, their heir taking
// over or their line ending.
func (f *farm) playerAges(day int) {
	if calendar.Gone(f.self.life, day) {
		f.passOn(day)
		return
	}
	if who := f.self.who(day); who.Stage != f.drawn {
		f.redressPlayer(who)
	}
}

// passOn is the player dying of old age on day: their eldest living child,
// or the eldest heir of their children, takes over, at home, with the
// money, the bag and the farm; with no one to, the line ends.
func (f *farm) passOn(day int) {
	old := f.self
	f.family.Die(old.id)
	id, ok := f.family.Heir(old.id)
	h := f.personByID(id)
	if !ok || h == nil {
		f.endLine(old)
		return
	}
	f.drop(h)
	h.ours = false
	f.self = &self{id: h.id, name: h.name, woman: h.woman, genes: h.who.Genes, life: h.life}
	f.forgetBonds()
	f.goHome()
	f.redressPlayer(f.self.who(day))
	f.notify(old.name + " morreu de velhice. " + h.name + " assume a fazenda.")
}

// personByID is the villager of id, or nil.
func (f *farm) personByID(id string) *villager {
	for _, v := range f.people {
		if v.id == id {
			return v
		}
	}
	return nil
}

// forgetBonds starts the new player's bonds afresh: the villagers knew
// the one before, not them.
func (f *farm) forgetBonds() {
	for _, v := range f.people {
		v.bond = bond.Bond{}
	}
}

// redressPlayer draws the player as who, and notes the time of life they
// are drawn at.
func (f *farm) redressPlayer(who Who) {
	f.drawn = who.Stage
	dress(f.player, f.art.Person(who))
}
