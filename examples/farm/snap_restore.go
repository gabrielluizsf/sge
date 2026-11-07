package farm

import (
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// restore puts snapshot s in farm f, which should be a new one. s is
// checked first, all of it, and an error leaves f as it was; so a hostile
// save gives an error, never a broken game.
func restore(f *farm, s *snapshot) error {
	l, err := check(f, s)
	if err != nil {
		return err
	}
	applyWorld(f, s, l)
	applyPeople(f, s, l)
	applyHouse(f, s, l)
	applyHerd(f, l)
	applyPlayer(f, s)
	return nil
}

// applyWorld sets the clock, the plots, the trees and the errand.
func applyWorld(f *farm, s *snapshot, l *loaded) {
	f.clock.Day, f.clock.Hour = s.Day, s.Hour
	f.seen, f.fainted, f.ids = s.Seen, s.Fainted, s.Ids
	*f.luck, *f.fate = random.Source(s.Luck), random.Source(s.Fate)
	f.plots, f.hits, f.family = l.plots, l.hits, l.tree
	for _, p := range f.plots {
		f.show(p)
	}
	f.stumps = l.stumps
	for _, st := range f.stumps {
		f.setThing(st.col, st.row, leftBy(st.tree))
	}
	f.errand, f.errandDone = orderOf(s.Errand), s.ErrandDone
	f.prog = progressOf(s)
	f.vit = vitalsOf(s)
	f.shown = f.xpValue()
	f.settled = map[*town.Town]bool{}
	for _, t := range l.towns {
		f.settled[t] = true
	}
	f.away = l.away
}

// applyPeople replaces the villagers, destroying those seen on screen.
func applyPeople(f *farm, s *snapshot, l *loaded) {
	for _, v := range f.people {
		if v.obj != nil {
			v.obj.Destroy()
		}
	}
	f.people = l.people
	f.notices = nil
}

// applyPlayer sets the player's kit, hotbar, place and who they are.
func applyPlayer(f *farm, s *snapshot) {
	*f.kit = *kitOf(s.Player)
	copy(f.pack.Hotbar, s.Hotbar)
	f.pack.Active = s.Active
	me := s.Self
	f.self = &self{id: me.ID, name: me.Name, woman: me.Woman, genes: me.Genes, life: lifeOf(me.Life)}
	if me.Spouse >= 0 {
		f.self.spouse = f.people[me.Spouse]
	}
	if me.Expecting != nil {
		f.self.expecting = &lineage.Expecting{Parents: me.Expecting.Parents, Due: me.Expecting.Due}
	}
	for _, v := range f.people {
		f.age(v, f.clock.Day)
	}
	f.talk, f.shop.open, f.question = nil, false, nil
	pos := f.player.Transform.Position
	f.player.Transform.Position = f.freed(vector.Vector3{X: s.At.X, Y: s.At.Y, Z: pos.Z})
	f.redressPlayer(f.self.who(f.clock.Day))
}
