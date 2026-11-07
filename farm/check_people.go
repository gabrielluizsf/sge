package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/random"
)

// knownRole reports whether role is a trade of the game.
func knownRole(role string) bool {
	_, ok := schedules[role]
	return ok
}

// checkPeople makes the villagers, each with the town and house they
// belong to, their partners and their ids in the family tree.
func (c *checker) checkPeople() error {
	n := len(c.s.People)
	if n > maxPeople {
		return fmt.Errorf("save: %d people", n)
	}
	ids := map[string]bool{}
	for _, e := range c.s.Family {
		ids[e.ID] = true
	}
	for i, p := range c.s.People {
		v, err := c.checkPerson(p)
		if err != nil {
			return fmt.Errorf("villager %d: %w", i, err)
		}
		if !ids[p.ID] {
			return fmt.Errorf("save: villager %d is not in the family tree", i)
		}
		c.out.people = append(c.out.people, v)
	}
	for i, p := range c.s.People {
		if err := c.link(i, p); err != nil {
			return err
		}
	}
	return nil
}

// link joins villager i to their partner, who has to have them as theirs.
func (c *checker) link(i int, p villagerState) error {
	if err := inRange("partner", p.Partner, -1, len(c.s.People)-1); err != nil {
		return err
	}
	if p.Partner >= 0 {
		if p.Partner == i || c.s.People[p.Partner].Partner != i {
			return fmt.Errorf("save: villager %d and their partner are not each other's", i)
		}
		c.out.people[i].partner = c.out.people[p.Partner]
	}
	return nil
}

// checkPerson is villager p made, or what is wrong with them.
func (c *checker) checkPerson(p villagerState) (*villager, error) {
	if err := firstError(
		plain("name", p.Name), plain("id", p.ID), checkGenes("villager", p.Genes),
		checkKit("villager", p.Kit), checkOrder("order", p.Order), checkLife("villager", p.Life),
		checkPoint("villager", p.At), checkExpecting("child", p.Expecting),
		finite("wait", p.Wait, -maxSeconds, maxSeconds), finite("homeFor", p.HomeFor, -maxSeconds, maxSeconds), finite("fuel", p.Fuel, -maxSeconds, maxSeconds),
		inRange("courting", p.Bond.Courted, int(bond.Strangers), int(bond.Married)),
		inRange("affection", p.Bond.Affection, 0, bond.Most),
		inRange("gifted", p.Bond.Gifted, 0, maxDay+1), inRange("talked", p.Bond.Talked, 0, maxDay+1),
		inRange("wedding", p.Bond.Wedding, 0, maxDay+bond.WeddingAfter),
	); err != nil {
		return nil, err
	}
	if !knownRole(p.Role) {
		return nil, fmt.Errorf("save: no trade %q", p.Role)
	}
	v := &villager{
		name: p.Name, woman: p.Woman, role: p.Role, who: Who{Woman: p.Woman, Role: p.Role, Genes: p.Genes},
		kit: kitOf(p.Kit), order: orderOf(p.Order), wait: p.Wait,
		pos:    c.f.ground.CellCenter(0, 0),
		inside: p.Inside, homeFor: p.HomeFor, fuel: p.Fuel, id: p.ID,
		life: lifeOf(p.Life), ours: p.Ours, luck: random.New(p.Luck),
		bond: bond.State{Affection: p.Bond.Affection, Courted: bond.Stage(p.Bond.Courted), Gifted: p.Bond.Gifted, Talked: p.Bond.Talked, Wedding: p.Bond.Wedding}.Bond(),
	}
	v.pos.X, v.pos.Y = p.At.X, p.At.Y
	v.mind = &mind{v: v}
	if p.Expecting != nil {
		e := lineage.Expecting{Parents: p.Expecting.Parents, Due: p.Expecting.Due}
		v.expecting = &e
	}
	return v, c.place(v, p)
}

// place gives v the town and house the snapshot says.
func (c *checker) place(v *villager, p villagerState) error {
	if p.Town == nil {
		return inRange("home", p.Home, -1, -1)
	}
	t, err := c.townAt(p.Town)
	if err != nil {
		return err
	}
	if err := inRange("home", p.Home, -1, len(t.Houses)-1); err != nil {
		return err
	}
	v.town = t
	if p.Home >= 0 {
		v.home = &t.Houses[p.Home]
	}
	return nil
}
