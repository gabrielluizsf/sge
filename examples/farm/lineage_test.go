package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
)

// Going to bed unmarried, the player just sleeps; married, they are asked
// only whether to have a child, yes or no, and the night passes either
// way.
func TestTheChildQuestionOnlyWhenMarried(t *testing.T) {
	_, f := newGame(idle{})
	v := single(t, f)
	f.clock.Hour = 21
	f.goToBed()
	if f.question != nil || f.clock.Hour != wakeAt {
		t.Fatalf("unmarried: asked %v, at %v", f.question, f.clock.Hour)
	}
	wed(t, f, v)
	f.clock.Hour = 12
	f.goToBed()
	if f.question != nil {
		t.Fatal("asked by day")
	}
	f.clock.Hour = 21
	f.goToBed()
	q := f.question
	if q == nil || len(q.ask) != 1 || q.ask[0] != "Devemos ter um filho?" || len(q.answers) != 2 || q.answers[0].text != "Sim" || q.answers[1].text != "Não" {
		t.Fatalf("asked %+v", q)
	}
	q.choice = 1 // no
	day := f.clock.Day
	f.answerQuestion()
	if f.self.expecting != nil || f.clock.Day != day+1 || f.question != nil {
		t.Fatalf("said no: expecting %v, day %d", f.self.expecting, f.clock.Day)
	}
}

// Said yes, the child is born when due, not before, a baby in the
// player's house and their child in the family tree; and grows up with
// the years.
func TestTheChildIsBornWhenDue(t *testing.T) {
	_, f := newGame(idle{})
	wed(t, f, single(t, f))
	f.clock.Hour = 21
	f.goToBed()
	f.answerQuestion()
	e := f.self.expecting
	if e == nil || e.Due != f.clock.Day-1+lineage.Gestation {
		t.Fatalf("expecting %+v on day %d", e, f.clock.Day)
	}
	f.goToBed()
	if f.question != nil {
		t.Error("asked again while expecting")
	}
	f.passDays(e.Due - f.clock.Day - 1)
	if len(f.family.Children(f.self.id)) != 0 {
		t.Fatal("born before due")
	}
	f.passDays(1)
	kids := f.family.Children(f.self.id)
	if len(kids) != 1 {
		t.Fatal("not born when due")
	}
	kid := f.personByID(kids[0])
	if kid == nil || !kid.ours || kid.stage() != life.Baby || f.self.expecting != nil || !f.childOf(kid) {
		t.Fatalf("the child: %+v", kid)
	}
	f.passDays(life.ChildFrom * daysPerYear)
	if kid.stage() != life.Child {
		t.Errorf("two years on, a %v", kid.stage())
	}
}

// When the player dies of old age their eldest child takes over, at home,
// with the money, the bag and the farm, their bonds fresh.
func TestTheHeirTakesOver(t *testing.T) {
	g, f := newGame(idle{})
	spouse := single(t, f)
	wed(t, f, spouse)
	first := haveChild(t, f)
	haveChild(t, f)
	f.money = 345
	f.bag.Add(wood, 7)
	f.self.life.Span = yearsOf(f.self.life, f.clock.Day) + 0.5/daysPerYear
	old := f.self.id
	f.passDays(1)
	g.Step(frame)
	if f.self.id != first.id || f.self.name != first.name || f.personByID(first.id) != nil {
		t.Fatalf("the player is %s, not the eldest child %s", f.self.id, first.id)
	}
	if f.money != 345 || f.bag.Count(wood) != 7 {
		t.Errorf("inherited $%d, %d wood", f.money, f.bag.Count(wood))
	}
	if f.room == nil || !f.room.own {
		t.Error("not at home")
	}
	if spouse.bond.MayHaveChild() || f.self.spouse != nil || !f.family.Kin(f.self.id, old) {
		t.Error("the heir took on their parent's marriage")
	}
	if len(f.notices) == 0 || f.notices[len(f.notices)-1].text != "Joaquim morreu de velhice. "+first.name+" assume a fazenda." {
		t.Errorf("notices %v", f.notices)
	}
}

// With no heir, the line ends: the player is asked only to start again,
// and a new farmer of twenty comes to a bare farm.
func TestTheLineEndsWithoutAnHeir(t *testing.T) {
	_, f := newGame(idle{})
	f.money = 999
	f.plots[0].crop = 0
	f.self.life.Span = yearsOf(f.self.life, f.clock.Day) + 0.5/daysPerYear
	old := f.self.id
	f.passDays(1)
	q := f.question
	if q == nil || q.ask[0] != "Fim da linhagem" || len(q.answers) != 1 {
		t.Fatalf("asked %+v", q)
	}
	f.answerQuestion()
	if f.self.id == old || f.self.who(f.clock.Day).Years != founderYears || f.money != 20 || f.plots[0].crop != -1 || f.question != nil {
		t.Errorf("restarted as %s of %d, $%d", f.self.id, f.self.who(f.clock.Day).Years, f.money)
	}
}

// F9 passes a year at once.
func TestAYearPassesAtOnce(t *testing.T) {
	_, f := newGame(idle{})
	day := f.clock.Day
	f.passDays(daysPerYear)
	if f.clock.Day != day+daysPerYear || f.seen != f.clock.Day {
		t.Errorf("day %d, lived to %d", f.clock.Day, f.seen)
	}
}
