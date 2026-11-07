package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// Every shed of a town stands apart from the others, solid, with its front
// open and reached on foot from the plaza.
func TestEveryShedCanBeWalkedTo(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	var last rect
	for i, name := range sheds {
		r := shedRect(f.art, tn, i)
		if i > 0 && r.col < last.col+last.cols {
			t.Errorf("%s at %+v overlaps %+v", name, r, last)
		}
		last = r
		if f.thingAt(r.col, r.row) != tileID["block"] {
			t.Errorf("%s is not built", name)
		}
		front := pathfind.Cell{Col: r.col + r.cols/2, Row: r.row + r.rows}
		if f.grid().Path(pathfind.Cell(tn.Center), front) == nil {
			t.Errorf("no way to the %s at %+v", name, front)
		}
		f.player.Transform.Position = f.center(front)
		if at, got := f.shedBy(f.player); at != tn || got != name {
			t.Errorf("before the %s: %v", name, got)
		}
		if words, _, ok := f.interaction(f.player); !ok || words == "" {
			t.Errorf("nothing to do before the %s", name)
		}
	}
}

// Every town has its baker, in white, and butcher, keeping their shops.
func TestATownHasItsShopkeepers(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	f.wake(tn)
	seen := map[string]*villager{}
	for _, v := range f.people {
		if v.town == tn {
			seen[v.role] = v
		}
	}
	for _, role := range []string{roleDealer, roleBaker, roleButcher} {
		v := seen[role]
		if v == nil {
			t.Fatalf("no %s", role)
		}
		if j := shopPost(f, v, 3); j == nil {
			t.Errorf("the %s has no post", role)
		}
	}
	if seen[roleBaker].who.Role != roleBaker {
		t.Error("the baker is not dressed as one")
	}
}

// A baker buys their wood, from a lumberjack, and does not fell trees.
func TestABakerBuysWood(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	f.wake(tn)
	var b, lj *villager
	for _, v := range f.people {
		if v.town == tn && v.role == roleBaker {
			b = v
		}
		if v.town == tn && v.role == roleLumberjack {
			lj = v
		}
	}
	if b == nil || lj == nil {
		t.Fatal("no baker or no lumberjack")
	}
	b.kit.money, b.fuel = 100, 0
	lj.kit.bag.Add(wood, 5)
	lj.inside = false
	lj.pos = f.center(pathfind.Cell(tn.Center))
	j := b.mind.bakerWood(f)
	if j == nil {
		t.Fatal("no job")
	}
	b.pos = lj.pos
	j.then(f)
	if b.kit.bag.Count(wood) != 5 || b.kit.money != 100-5*woodPrice {
		t.Errorf("%d wood, $%d", b.kit.bag.Count(wood), b.kit.money)
	}
	j = b.mind.bakerWood(f)
	j.then(f)
	if b.fuel != 5 {
		t.Errorf("fuel %v", b.fuel)
	}
	money := b.kit.money
	b.mind.bakerWood(f).then(f)
	if b.fuel != 4 || b.kit.money != money+breadPay {
		t.Errorf("fuel %v money %d", b.fuel, b.kit.money)
	}
}

// Wheat bakes into bread at a forge, which feeds and sells for more; the
// bakery buys wheat dearer than a stall, and sells bread.
func TestBread(t *testing.T) {
	_, f := newGame(idle{})
	tn := f.towns.Of(1, 0)
	c, _ := cropIndex("wheat")
	f.bag.Add(cropOf(c), 2)
	f.player.Transform.Position = f.center(pathfind.Cell{Col: tn.Forge.Col, Row: tn.Forge.Row + 1})
	words, do, ok := f.forgeNear(f.player)
	if !ok || words == "" {
		t.Fatal("no baking at the forge")
	}
	do()
	if f.bag.Count(bread) != 1 || f.bag.Count(cropOf(c)) != 1 {
		t.Fatalf("%d bread, %d wheat", f.bag.Count(bread), f.bag.Count(cropOf(c)))
	}
	g, _ := goodOf(bread)
	if g.price <= crops[c].sellPrice || g.nutrition <= crops[c].nutrition {
		t.Errorf("bread is worth no more than wheat: %+v", g)
	}
	f.kit.money = 0
	_, lines := f.shedLines("bakery")
	if !lines[1].do(f, f.kit) || f.kit.money != crops[c].sellPrice+wheatPrem {
		t.Errorf("money %d", f.kit.money)
	}
	if lines[0].do(f, f.kit) && f.kit.money < 0 {
		t.Error("in debt")
	}
	f.kit.money = 100
	before := f.bag.Count(bread)
	if !lines[0].do(f, f.kit) || f.bag.Count(bread) != before+1 || f.kit.money != 100-g.cost {
		t.Errorf("%d bread, money %d", f.bag.Count(bread), f.kit.money)
	}
}

// Bread in hand is eaten as a crop is, and it is what is kept in a save; a
// save holding a thing not to eat is refused.
func TestBreadIsEaten(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add(bread, 1)
	f.bag.Add(wool, 1)
	f.pack.Hotbar[4] = wool
	if f.hold(4) {
		t.Error("wool in hand to eat")
	}
	f.pack.Hotbar[4] = bread
	if !f.hold(4) || f.slot != slotFood || f.food != bread {
		t.Fatalf("slot %d food %q", f.slot, f.food)
	}
	if s := takeSnapshot(f); s.Player.Food != bread {
		t.Errorf("saved %q", s.Player.Food)
	}
	s := takeSnapshot(f)
	s.Player.Food = "wool"
	_, g := newGame(idle{})
	if err := restore(g, &s); err == nil {
		t.Error("wool was let in as food")
	}
	f.vit.Food = 10
	col, row := f.cellOf(feet(f.player))
	if f.useAs(f.kit, col, row) != actEat || f.vit.Food != 50 || f.bag.Count(bread) != 0 {
		t.Errorf("food %v, bread %d", f.vit.Food, f.bag.Count(bread))
	}
}

// One's own animal sells to the butcher, or is killed for its meat with
// blows; a horse gives none, and the butcher sells meat.
func TestTheButcher(t *testing.T) {
	_, f := newGame(idle{})
	c := f.herd.crits[0]
	f.kit.money = 0
	if !f.sellAnimal(f.kit, "cow") || f.kit.money != animalPrice["cow"] || len(f.herd.crits) != 0 {
		t.Fatalf("money %d, %d animals", f.kit.money, len(f.herd.crits))
	}
	if f.sellAnimal(f.kit, "cow") {
		t.Error("sold a cow that is gone")
	}
	f.herd.crits = append(f.herd.crits, c)
	c.obj = nil
	at := c.pos.Add(vector.Vector3{X: -tile})
	for range 10 {
		if len(f.herd.crits) == 0 {
			break
		}
		f.strikeAnimal(at, vector.Vector3{X: 1}, axeHarm)
	}
	if len(f.herd.crits) != 0 || f.bag.Count(beef) != meatOf["cow"].n {
		t.Errorf("%d animals, %d beef", len(f.herd.crits), f.bag.Count(beef))
	}
	h := f.herd.add("horse", 0, f.herd.pens[0])
	if f.strikeAnimal(h.pos.Add(vector.Vector3{X: -tile}), vector.Vector3{X: 1}, axeHarm) {
		t.Error("a blow landed on the horse")
	}
	_, lines := f.shedLines("butcher")
	f.kit.money = 100
	if !lines[0].do(f, f.kit) || f.bag.Count(beef) != meatOf["cow"].n+1 {
		t.Errorf("no beef bought: %d", f.bag.Count(beef))
	}
}
