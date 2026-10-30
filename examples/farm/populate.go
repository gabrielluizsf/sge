package farm

import (
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// populate gives town t's people, just settled, their lives: each an age
// their trade allows, drawn by the town's seed; puts them in the family
// tree, the grown-ups first; and pairs some of the grown-ups as couples,
// the young of the town their first couple's children.
func (f *farm) populate(t *town.Town, folk []*villager) {
	luck := random.New(t.Seed + 13)
	day := f.clock.Day
	for i, v := range folk {
		v.id = "t" + itoa(int(t.Seed)) + "-" + itoa(i)
		v.life = calendar.Draw(luck, drawYears(luck, v.role), day)
		f.age(v, day)
	}
	pairUp(folk, luck)
	var parents [2]string
	for _, v := range folk {
		if v.stage() >= life.Adult {
			f.family.Add(lineage.Person{ID: v.id, Born: v.life.Born})
			if v.partner != nil && parents[0] == "" && v.woman {
				parents = [2]string{v.id, v.partner.id}
			}
		}
	}
	for _, v := range folk {
		if v.stage() < life.Adult {
			f.family.Add(lineage.Person{ID: v.id, Parents: parents, Born: v.life.Born})
		}
	}
}

// drawYears is how old someone of role is when their town is first
// settled: the smith and lumberjacks of working age, farmers grown, some
// old, and the wanderers of any age, children and elders among them.
func drawYears(luck *random.Source, role string) float64 {
	u := luck.Float()
	switch role {
	case roleSmith:
		return 25 + 30*u
	case roleLumberjack:
		return 20 + 35*u
	case roleFarmer:
		return 18 + 60*u
	}
	switch w := luck.Float(); {
	case w < 0.3:
		return life.ChildFrom + 2 + 8*u
	case w < 0.4:
		return life.TeenFrom + 6*u
	case w < 0.75:
		return life.AdultFrom + 40*u
	}
	return life.ElderFrom + 25*u
}

// pairUp makes couples of the grown-up men and women of folk, about half
// of those who could be, as luck falls.
func pairUp(folk []*villager, luck *random.Source) {
	var men []*villager
	for _, v := range folk {
		if !v.woman && v.stage() >= life.Adult {
			men = append(men, v)
		}
	}
	for _, w := range folk {
		if !w.woman || w.stage() < life.Adult || len(men) == 0 {
			continue
		}
		if luck.Float() < 0.5 {
			w.partner, men[0].partner = men[0], w
			men = men[1:]
		}
	}
}
