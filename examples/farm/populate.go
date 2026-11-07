package farm

import (
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// populate gives town t's people, just settled, their lives: each household's
// couple grown, the woman young enough to have borne the child, and the
// child young; puts them in the family tree, the grown-ups first, the
// child with the couple as parents.
func (f *farm) populate(t *town.Town, homes []household) {
	luck := random.New(t.Seed + 13)
	day := f.clock.Day
	for i, h := range homes {
		for j, v := range []*villager{h.man, h.woman, h.child} {
			v.id = "t" + itoa(int(t.Seed)) + "-" + itoa(3*i+j)
			v.life = calendar.Draw(luck, working(v, drawYears(luck, j)), day)
			f.age(v, day)
		}
		h.man.partner, h.woman.partner = h.woman, h.man
	}
	for _, h := range homes {
		f.family.Add(lineage.Person{ID: h.man.id, Born: h.man.life.Born})
		f.family.Add(lineage.Person{ID: h.woman.id, Born: h.woman.life.Born})
	}
	for _, h := range homes {
		f.family.Add(lineage.Person{ID: h.child.id, Parents: [2]string{h.woman.id, h.man.id}, Born: h.child.life.Born})
	}
}

// drawYears is how old member j of a household (0 the man, 1 the woman, 2
// the child) is when their town is first settled: the couple grown, the
// child of any age from a toddler to twenty, so some are grown and free to
// court the player, the child at the least sixteen years younger than the
// woman.
func drawYears(luck *random.Source, j int) float64 {
	u := luck.Float()
	switch j {
	case 0:
		return 38 + 30*u
	case 1:
		return 36 + 24*u
	}
	return life.ChildFrom + 1 + (20-life.ChildFrom-1)*u
}

// working is years for v, less than an elder's when they are the smith, who
// works the anvil.
func working(v *villager, years float64) float64 {
	if v.smith() {
		return min(years, life.ElderFrom-2)
	}
	return years
}
