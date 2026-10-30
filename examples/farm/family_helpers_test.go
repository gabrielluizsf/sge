package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/life"
)

// single is a grown-up villager of the towns near the farm living as no
// one's partner, the player not kin to them.
func single(t *testing.T, f *farm) *villager {
	t.Helper()
	for c := -1; c <= 1; c++ {
		if tn := f.towns.Of(c, 0); tn != nil {
			f.settle(tn)
		}
	}
	for _, v := range f.people {
		if v.partner == nil && !v.ours && v.stage() == life.Adult && !v.smith() {
			return v
		}
	}
	t.Fatal("no single grown-up in the towns")
	return nil
}

// wed marries the player to v, as fond of them as can be: asked out,
// proposed to, and the wedding days passed.
func wed(t *testing.T, f *farm, v *villager) {
	t.Helper()
	v.bond.Affection = bond.Most
	if err := f.askOut(v); err != nil {
		t.Fatalf("asked out: %v", err)
	}
	if err := f.propose(v); err != nil {
		t.Fatalf("proposed: %v", err)
	}
	f.passDays(bond.WeddingAfter)
	if f.self.spouse != v || !v.ours {
		t.Fatalf("not wed: spouse %v", f.self.spouse)
	}
}

// haveChild has the married player decide on a child tonight and passes
// the days until they are born.
func haveChild(t *testing.T, f *farm) *villager {
	t.Helper()
	f.clock.Hour = 21
	f.goToBed()
	if f.question == nil {
		t.Fatal("not asked whether to have a child")
	}
	f.answerQuestion() // the first answer, yes
	n := len(f.family.Children(f.self.id))
	f.passDays(f.self.expecting.Due - f.clock.Day)
	kids := f.family.Children(f.self.id)
	if len(kids) <= n {
		t.Fatal("no child born")
	}
	return f.personByID(kids[len(kids)-1])
}
