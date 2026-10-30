package main

import (
	"testing"

	"github.com/gabrielluizsf/sge/inventory"
	"github.com/gabrielluizsf/sge/request"
	"github.com/gabrielluizsf/sge/town"
)

// eastTown is the town east of the farm, settled, with its people.
func eastTown(f *farm) *town.Town {
	t := f.towns.Of(1, 0)
	f.settle(t)
	return t
}

// tradesman is town t's first villager of the trade.
func tradesman(f *farm, t *town.Town, role string) *villager {
	for _, v := range f.people {
		if v.town == t && v.role == role {
			return v
		}
	}
	return nil
}

// Every town has a smith and a lumberjack, women and men of many skins,
// and a garden for each house.
func TestTownsArePeopled(t *testing.T) {
	_, f := newGame(idle{})
	tn := eastTown(f)
	if tradesman(f, tn, roleSmith) == nil || tradesman(f, tn, roleLumberjack) == nil {
		t.Fatal("no smith or lumberjack")
	}
	women, skins := 0, map[string]bool{}
	for r := range 4 {
		for c := -3; c < 4; c++ {
			if tn := f.towns.Of(c, r); tn != nil {
				f.settle(tn)
			}
		}
	}
	for _, v := range f.people {
		if v.woman {
			women++
		}
		skins[v.look.Skin[0]] = true
	}
	if women == 0 || women == len(f.people) || len(skins) < 4 {
		t.Errorf("%d women of %d, %d skins", women, len(f.people), len(skins))
	}
	for _, h := range tn.Houses {
		if f.plotAt(h.Garden.Col, h.Garden.Row) == nil {
			t.Errorf("no garden at %v", h.Garden)
		}
	}
}

// A villager's order is delivered for its reward, whole or not at all,
// and they ask again later.
func TestAnOrderIsDelivered(t *testing.T) {
	g, f := newGame(idle{})
	v := tradesman(f, eastTown(f), roleFarmer)
	if v == nil {
		v = tradesman(f, eastTown(f), roleWanderer)
	}
	v.order = request.Order{Want: inventory.List{{Item: wood, N: 2}}, Reward: 12}
	f.talk = &talk{v: v}
	f.deliver(v)
	if f.money != 20 || len(v.order.Want) == 0 {
		t.Fatal("delivered without the wood")
	}
	f.bag.Add(wood, 3)
	f.deliver(v)
	if f.money != 32 || f.bag.Count(wood) != 1 || len(v.order.Want) != 0 || v.wants(f) {
		t.Fatalf("money %d, wood %d, order %+v", f.money, f.bag.Count(wood), v.order)
	}
	for range int(2*orderPause/frame) + 2 {
		g.Step(frame)
	}
	if len(v.order.Want) == 0 {
		t.Error("they never asked again")
	}
}

// The smith gives the stone axe for their errand, then sells the better
// axes, each only once and only for its price.
func TestTheSmithsAxes(t *testing.T) {
	_, f := newGame(idle{})
	s := tradesman(f, eastTown(f), roleSmith)
	f.talk = &talk{v: s}
	f.doErrand(s)
	if f.axe != noAxe || f.errandDone {
		t.Fatal("the errand was done without its items")
	}
	for _, w := range f.errand.Want {
		f.bag.Add(w.Item, w.N)
	}
	f.doErrand(s)
	if f.axe != 0 || !f.errandDone || s.wants(f) {
		t.Fatalf("axe %d after the errand", f.axe)
	}
	f.money = 100
	f.buyAxe(2, s.pos) // gold, $150
	f.buyAxe(1, s.pos) // iron, $60
	f.buyAxe(1, s.pos)
	if f.axe != 1 || f.money != 40 {
		t.Errorf("axe %d, money %d", f.axe, f.money)
	}
}

// The player trades wood with a lumberjack: buying dear, selling cheap,
// as long as one has the wood and the other the money.
func TestWoodIsTradedWithTheLumberjack(t *testing.T) {
	_, f := newGame(idle{})
	lj := tradesman(f, eastTown(f), roleLumberjack)
	lj.kit.bag.Add(wood, 1)
	f.tradeWood(deal{lj.kit, f.kit, lumberSells}, lj.pos)
	f.tradeWood(deal{lj.kit, f.kit, lumberSells}, lj.pos)
	if f.bag.Count(wood) != 1 || f.money != 20-lumberSells || lj.kit.bag.Count(wood) != 0 {
		t.Fatalf("bought: wood %d, money %d", f.bag.Count(wood), f.money)
	}
	lj.kit.money = 0
	f.tradeWood(deal{f.kit, lj.kit, lumberBuys}, lj.pos)
	if f.bag.Count(wood) != 1 {
		t.Fatal("sold to one with no money")
	}
	lj.kit.money = 10
	f.tradeWood(deal{f.kit, lj.kit, lumberBuys}, lj.pos)
	if f.bag.Count(wood) != 0 || lj.kit.bag.Count(wood) != 1 || f.money != 20-lumberSells+lumberBuys {
		t.Errorf("sold: wood %d, money %d", f.bag.Count(wood), f.money)
	}
}

// Left to their days, the people of the town work: lumberjacks fell
// trees, the smith keeps the forge going, and at night all are at home.
func TestTheTownLivesItsDay(t *testing.T) {
	g, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	g.Find("farmer").Transform.Position = standIn(f, tn.Center.Col, tn.Center.Row+1)
	f.clock.Hour = 8
	for range 60 * 120 {
		g.Step(frame)
	}
	felled, forged := len(f.stumps), 0
	for _, v := range f.people {
		if v.town == tn && v.role == roleSmith && v.kit.money != 40 {
			forged++
		}
	}
	if felled == 0 || forged == 0 {
		t.Errorf("in two minutes: %d trees felled, smith at work %d", felled, forged)
	}
	f.clock.Hour = 22.5
	for range 60 * 60 {
		g.Step(frame)
	}
	for _, v := range f.people {
		if v.town == tn && v.home != nil && !v.inside {
			t.Errorf("%s the %s is out at night, at %v", v.name, v.role, v.pos)
		}
	}
}
