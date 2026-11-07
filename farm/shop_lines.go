package farm

import (
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/vector"
)

// shopLine is a line of a shed's menu: its text, and what choosing it does.
type shopLine struct {
	text string
	do   func(f *farm, k *kit) bool
}

// wheatPrem is how much more the bakery pays for wheat than a stall.
const wheatPrem = 4

// animalPrice is what a butcher pays for an animal of the player's.
var animalPrice = map[string]int{"cow": 240, "sheep": 90}

// meatOf is the meat and how much an animal of a kind gives, killed.
var meatOf = map[string]struct {
	item string
	n    int
}{"cow": {beef, 3}, "sheep": {mutton, 2}}

// shedLines are the title and the lines of the menu of the shed kind.
func (f *farm) shedLines(kind string) (string, []shopLine) {
	var lines []shopLine
	switch kind {
	case "stable":
		for _, o := range stableOffers {
			lines = append(lines, shopLine{i18n.T("%s $%d", i18n.T(o.label), o.price), func(f *farm, k *kit) bool { return f.buyAnimal(k, o) }})
		}
		return i18n.T("Stable"), lines
	case "bakery":
		c, _ := cropIndex("wheat")
		price := crops[c].sellPrice + wheatPrem
		return i18n.T("Bakery"), []shopLine{
			{i18n.T("Buy %s $%d", i18n.T("Bread"), goodCost(bread)), buys(bread)},
			{i18n.T("Sell wheat $%d each (you have %d)", price, f.bag.Count(cropOf(c))), func(f *farm, k *kit) bool { return f.sellWheat(k, price) }},
		}
	}
	lines = []shopLine{
		{i18n.T("Buy %s $%d", i18n.T("Beef"), goodCost(beef)), buys(beef)},
		{i18n.T("Buy %s $%d", i18n.T("Mutton"), goodCost(mutton)), buys(mutton)},
	}
	for _, kind := range []string{"cow", "sheep"} {
		lines = append(lines, shopLine{i18n.T("Sell your %s $%d", i18n.T(animalLabel[kind]), animalPrice[kind]), func(f *farm, k *kit) bool { return f.sellAnimal(k, kind) }})
	}
	return i18n.T("Butcher"), lines
}

// animalLabel are the animals' names on screen, as English keys.
var animalLabel = map[string]string{"cow": "Cow", "sheep": "Sheep", "horse": "Horse"}

// goodCost is what item costs at a shop.
func goodCost(item string) int {
	g, _ := goodOf(item)
	return g.cost
}

// buys is a line's buying one of item.
func buys(item string) func(f *farm, k *kit) bool {
	return func(f *farm, k *kit) bool {
		at := f.player.Transform.Position
		at.Y -= tile
		cost := goodCost(item)
		if k.money < cost {
			f.float(at, i18n.T("Not enough money"), 'r')
			f.sound("deny", nil)
			return false
		}
		k.money -= cost
		k.bag.Add(item, 1)
		f.float(at, i18n.T("-$%d", cost), 'r')
		f.sound("coin", &at)
		return true
	}
}

// sellWheat sells k's wheat to the bakery, at price each.
func (f *farm) sellWheat(k *kit, price int) bool {
	c, _ := cropIndex("wheat")
	n := k.bag.Count(cropOf(c))
	at := f.player.Transform.Position
	at.Y -= tile
	if n == 0 {
		f.float(at, i18n.T("Nothing to sell"), 'r')
		f.sound("deny", nil)
		return false
	}
	k.bag.Take(cropOf(c), n)
	k.sold += n
	k.money += n * price
	f.earnBy(k, xpSale(n*price))
	f.float(at, i18n.T("+$%d", n*price), 'y')
	f.sound("sell", &at)
	return true
}

// sellAnimal sells the butcher one of the player's animals of a kind.
func (f *farm) sellAnimal(k *kit, kind string) bool {
	at := f.player.Transform.Position
	at.Y -= tile
	c := f.herd.own(kind)
	if c == nil {
		f.float(at, i18n.T("You have none"), 'r')
		f.sound("deny", nil)
		return false
	}
	f.herd.remove(c)
	k.money += animalPrice[kind]
	f.float(at, i18n.T("+$%d", animalPrice[kind]), 'y')
	f.sound("sell", &at)
	return true
}

// own is one of the player's animals of kind, or nil.
func (h *herd) own(kind string) *creature {
	for _, c := range h.crits {
		if c.kind == kind && c.pen.town == nil && !c.held {
			return c
		}
	}
	return nil
}

// remove takes c out of the herd, for good.
func (h *herd) remove(c *creature) {
	if c.obj != nil {
		c.obj.Destroy()
	}
	for i, o := range h.crits {
		if o == c {
			h.crits = append(h.crits[:i], h.crits[i+1:]...)
			return
		}
	}
}

// animalHurt is how many hits of harm one of the player's animals takes
// before it falls.
const animalHurt = 30.0

// strikeAnimal is the player's blow at their own animal ahead of them, if
// there is one that gives meat: it bleeds, and in the end falls and leaves
// its meat. It reports whether the blow landed.
func (f *farm) strikeAnimal(at, facing vector.Vector3, harm float64) bool {
	if f.herd == nil {
		return false
	}
	for _, c := range f.herd.crits {
		d := c.pos.Sub(at)
		l := d.Length()
		if _, meat := meatOf[c.kind]; !meat || c.pen.town != nil || c.held || l > 1.6*tile || l > 0 && d.Dot(facing)/l < 0.3 {
			continue
		}
		f.bleed(c.pos, at, 14)
		if c.hurt += harm; c.hurt >= animalHurt {
			m := meatOf[c.kind]
			f.bleed(c.pos, at, 30)
			f.bag.Add(m.item, m.n)
			g, _ := goodOf(m.item)
			f.notify(i18n.T("+%d %s", m.n, i18n.T(g.label)))
			f.herd.remove(c)
		}
		return true
	}
	return false
}

// shedWords are what E says before each shed, as English keys.
var shedWords = map[string]string{"stable": "E: buy animals", "bakery": "E: enter the bakery", "butcher": "E: buy meat"}
