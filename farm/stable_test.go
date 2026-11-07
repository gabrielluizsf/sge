package farm

import (
	"encoding/json"
	"testing"

	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// Every town's stable and pen are built, solid, and reached from its plaza
// on foot, in both looks' layouts.
func TestAStableCanBeWalkedTo(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	r := stableRect(f.art, tn)
	if r.cols == 0 || f.thingAt(r.col, r.row) != tileID["block"] {
		t.Fatalf("no stable at %+v", r)
	}
	p := townPen(tn)
	if !isFence(f.thingAt(p.col, p.row)) || isFence(f.thingAt(p.col+p.cols/2, p.row+p.rows-1)) {
		t.Errorf("the pen's fence: %d", f.thingAt(p.col, p.row))
	}
	front := pathfind.Cell{Col: r.col, Row: r.row + r.rows}
	if !f.open(front.Col, front.Row) {
		t.Fatalf("the stable's front is shut: %+v", front)
	}
	if way := f.grid().Path(pathfind.Cell(tn.Center), front); way == nil {
		t.Errorf("no way from the plaza to the stable at %+v", front)
	}
	if way := f.grid().Path(pathfind.Cell{Col: front.Col, Row: front.Row}, pathfind.Cell{Col: p.col + p.cols/2, Row: p.row + 1}); way == nil {
		t.Error("no way through the pen's gate")
	}
}

// A town's stable is kept by a dealer and has animals in its pen, which are
// not the player's to milk or save.
func TestAStableHasADealerAndAnimals(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	f.wake(tn)
	dealers := 0
	for _, v := range f.people {
		if v.town == tn && v.role == roleDealer {
			dealers++
		}
	}
	if dealers != 1 {
		t.Errorf("%d dealers", dealers)
	}
	herdSize := func() int {
		n := 0
		for _, c := range f.herd.crits {
			if c.pen.town == tn {
				n++
			}
		}
		return n
	}
	if herdSize() != 3 {
		t.Errorf("%d animals at the stable", herdSize())
	}
	f.wake(tn)
	if herdSize() != 3 {
		t.Error("stocked twice")
	}
	if s := takeSnapshot(f); len(s.Herd) != 1 || len(s.Pens) != 1 {
		t.Errorf("saved %d animals in %d pens", len(s.Herd), len(s.Pens))
	}
}

// Before the stable, E offers the animals; buying takes the money and puts
// the animal in a pen of its kind on the farm, the next in the same.
func TestBuyingAnAnimal(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	r := stableRect(f.art, tn)
	f.player.Transform.Position = f.center(pathfind.Cell{Col: r.col, Row: r.row + r.rows})
	if !atShed(f, tn, "stable") {
		t.Fatal("not at the stable")
	}
	words, do, ok := f.interaction(f.player)
	if !ok || words == "" {
		t.Fatal("nothing offered")
	}
	do()
	if !f.shop.open || f.shop.kind != "stable" {
		t.Fatal("no menu")
	}
	f.kit.money = 10
	if f.buyAnimal(f.kit, stableOffers[0]) || len(f.herd.crits) != 1 {
		t.Fatal("bought with no money")
	}
	f.kit.money = 1000
	if !f.buyAnimal(f.kit, stableOffers[0]) || !f.buyAnimal(f.kit, stableOffers[0]) || f.kit.money != 1000-2*stableOffers[0].price {
		t.Fatalf("money %d", f.kit.money)
	}
	a, b := f.herd.crits[1], f.herd.crits[2]
	if a.kind != "sheep" || a.pen != b.pen || a.pen.town != nil {
		t.Errorf("%+v %+v", a, b)
	}
	if !f.penFitsItsFence(a.pen) {
		t.Error("no fence round the new pen")
	}
}

// penFitsItsFence reports whether p's fence is up.
func (f *farm) penFitsItsFence(p *pen) bool {
	return isFence(f.thingAt(p.r.col, p.r.row))
}

// A pen is moved with the arranging: its fence and its animals go with it,
// the old ground left clear; on a tree or the house it may not stand.
func TestAPenIsMovedWithTheEditor(t *testing.T) {
	_, f := newGame(idle{})
	first := f.herd.pens[0]
	c := f.herd.crits[0]
	old := first.r
	f.player.Transform.Position = f.center(pathfind.Cell{Col: old.col + 2, Row: old.row + 2})
	if !f.startEdit("") {
		t.Fatal("no arranging")
	}
	f.edit.col, f.edit.row = old.col, old.row
	f.editAct()
	if f.edit.pen != first {
		t.Fatal("not lifted")
	}
	to, ok := f.freePenSpot()
	if !ok {
		t.Fatal("no free ground")
	}
	f.edit.col, f.edit.row = to.col, to.row
	if _, fits := f.penGhost(); !fits {
		t.Fatalf("no fit at %+v", to)
	}
	before := c.pos
	f.editAct()
	if first.r != to || f.edit.pen != nil {
		t.Fatalf("at %+v", first.r)
	}
	if isFence(f.thingAt(old.col, old.row)) || !isFence(f.thingAt(to.col, to.row)) {
		t.Error("the fence did not move")
	}
	if want := before.X + float64(to.col-old.col)*tile; c.pos.X != want {
		t.Errorf("the cow at %v, not %v", c.pos.X, want)
	}
	house := siteFor(f.art).home
	f.edit.col, f.edit.row = house.col, house.row
	f.edit.pen, f.edit.penCol, f.edit.penRow = first, 0, 0
	if _, fits := f.penGhost(); fits {
		t.Error("a pen over the house")
	}
}

// Pens moved and animals bought are kept in a save, and a pen on a house in
// one is refused.
func TestPensAreSaved(t *testing.T) {
	_, f := newGame(idle{})
	f.kit.money = 1000
	f.buyAnimal(f.kit, stableOffers[1])
	pen := f.herd.crits[1].pen.r
	data, err := json.Marshal(takeSnapshot(f))
	if err != nil {
		t.Fatal(err)
	}
	var s snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil {
		t.Fatal(err)
	}
	if len(g.herd.crits) != 2 || g.herd.crits[1].kind != "horse" || g.herd.crits[1].pen.r != pen || !isFence(g.thingAt(pen.col, pen.row)) {
		t.Errorf("%+v", g.herd.crits)
	}
	house := siteFor(f.art).home
	s.Pens[0].Col, s.Pens[0].Row = house.col, house.row
	_, h := newGame(idle{})
	if err := restore(h, &s); err == nil {
		t.Error("a pen on the house was let in")
	}
}

// Wool comes from one's own sheep, and grows back in days.
func TestShearing(t *testing.T) {
	_, f := newGame(idle{})
	f.kit.money = 1000
	f.buyAnimal(f.kit, stableOffers[0])
	c := f.herd.crits[1]
	f.player.Transform.Position.X, f.player.Transform.Position.Y = c.pos.X+tile, c.pos.Y
	if f.sheepBy(f.player) != c {
		t.Fatal("no sheep by")
	}
	if !f.shear(c) || f.bag.Count(wool) != 1 {
		t.Fatal("no wool")
	}
	if f.shear(c) || f.woolIn(c) != woolDays {
		t.Errorf("wool in %d days", f.woolIn(c))
	}
	f.clock.Day += woolDays
	if !f.shear(c) || f.bag.Count(wool) != 2 {
		t.Error("the wool did not grow back")
	}
}

// Goods sell at the stall, and the old world's save without a herd loads.
func TestGoodsSell(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(milk, 2)
	f.bag.Add(wool, 1)
	if got := worth(f.kit); got != 2*18+12 {
		t.Errorf("worth %d", got)
	}
	f.kit.sellOut()
	if f.bag.Count(milk) != 0 || f.kit.money < 2*18+12 {
		t.Errorf("milk %d money %d", f.bag.Count(milk), f.kit.money)
	}
}

// A save from before the herd, with a stump on ground since taken, loads.
func TestAnOldSaveLoads(t *testing.T) {
	_, f := newGame(idle{})
	s := takeSnapshot(f)
	s.Herd, s.Pens = nil, nil
	r := townPen(f.towns.Of(1, 0))
	s.Stumps = append(s.Stumps, stumpState{Col: r.col + 1, Row: r.row + 1, Tree: 0, Left: 1})
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil {
		t.Fatal(err)
	}
	if len(g.herd.crits) != 1 {
		t.Error("the first cow is gone")
	}
}

// atShed reports whether the player stands before the shed name of t.
func atShed(f *farm, t *town.Town, name string) bool {
	at, got := f.shedBy(f.player)
	return at == t && got == name
}
