package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/inventory"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/request"
	"github.com/gabrielluizsf/sge/vector"
)

// checker checks a snapshot against a farm, making what it holds as it
// goes, without changing the farm.
type checker struct {
	f   *farm
	s   *snapshot
	out loaded
}

// check is what snapshot s makes for farm f, or why s is not a save the
// game could have made.
func check(f *farm, s *snapshot) (*loaded, error) {
	c := &checker{f: f, s: s}
	steps := []func() error{
		func() error { return checkClock(s) },
		func() error { return checkProgress(s) },
		func() error { return checkVitals(s) },
		func() error { return checkKit("player", s.Player) },
		func() error { return checkHotbar(s) },
		func() error { return checkPoint("player", s.At) },
		func() error { return checkRoom(s) },
		c.checkFamily, c.checkTowns, c.checkAway, c.checkPeople,
		func() error { return checkSelf(s) },
		c.checkPlots, c.checkStumps, c.checkHits,
		func() error { return checkOrder("errand", s.Errand) },
		c.checkPlace, c.checkHouse, c.checkHerd,
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return nil, err
		}
	}
	return &c.out, nil
}

// checkFamily makes the family tree, with the player in it.
func (c *checker) checkFamily() error {
	if len(c.s.Family) > maxFamily {
		return fmt.Errorf("save: %d in the family tree", len(c.s.Family))
	}
	es := make([]lineage.Entry, 0, len(c.s.Family))
	for _, k := range c.s.Family {
		if err := firstError(plain("family id", k.ID), inRange("born", k.Born, -maxDay, maxDay)); err != nil {
			return err
		}
		es = append(es, lineage.Entry{Person: lineage.Person{ID: k.ID, Parents: k.Parents, Born: k.Born}, Dead: k.Dead})
	}
	t, err := lineage.Restore(es)
	if err != nil {
		return err
	}
	c.out.tree = t
	for _, k := range c.s.Family {
		if k.ID == c.s.Self.ID {
			return nil
		}
	}
	return fmt.Errorf("save: the player is not in the family tree")
}

// checkPlace is nil if the player stands where one may.
func (c *checker) checkPlace() error {
	col, row := c.f.cellOf(vec(c.s.At))
	// A fence may have been built since where the player was saved: the
	// player is moved out of it when loaded.
	if thing := c.f.thingAt(col, row); !isFence(thing) && c.f.tiles.IsSolid(thing) {
		return fmt.Errorf("save: the player is inside something solid")
	}
	return nil
}

// kitOf is the kit k stands for.
func kitOf(k kitState) *kit {
	out := &kit{money: k.Money, slot: k.Slot, seed: k.Seed, axe: k.Axe, sold: k.Sold, used: k.Used, exact: k.Exact}
	out.food = k.Food
	out.bag = bagOf(k.Bag)
	return out
}

// bagOf is a bag of l.
func bagOf(l []stackState) inventory.Bag {
	var b inventory.Bag
	for _, s := range l {
		b.Add(s.Item, s.N)
	}
	return b
}

// orderOf is the order o stands for.
func orderOf(o orderState) request.Order {
	var want inventory.List
	for _, s := range o.Want {
		want = append(want, inventory.Stack{Item: s.Item, N: s.N})
	}
	return request.Order{Want: want, Reward: o.Reward}
}

// lifeOf is the life l stands for.
func lifeOf(l lifeState) life.Life { return life.Life{Born: l.Born, Span: l.Span} }

// vec is the point p as a vector.
func vec(p point) vector.Vector3 { return vector.Vector3{X: p.X, Y: p.Y} }
