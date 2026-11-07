package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/i18n"
)

// craftBook is the Craft tab of the inventory over what the player
// carries: it lists the recipes against the bag and makes them from it.
type craftBook struct{ f *farm }

// Recipes are the recipes, each with what the bag has of its materials.
func (c craftBook) Recipes() []backpack.Recipe {
	out := make([]backpack.Recipe, 0, len(recipes))
	for _, r := range recipes {
		it := c.f.itemOf(r.item)
		rc := backpack.Recipe{ID: r.item, Name: it.Name, Icon: it.Icon, Info: it.Info}
		for _, n := range r.needs {
			m := c.f.itemOf(n.Item)
			rc.Needs = append(rc.Needs, backpack.Need{Name: m.Name, Icon: m.Icon, Want: n.N, Have: c.f.bag.Count(n.Item)})
		}
		out = append(out, rc)
	}
	return out
}

// Make takes the materials of the recipe for item out of the bag and puts
// the item in, if every material is there.
func (c craftBook) Make(item string) bool {
	r, ok := recipeOf(item)
	if !ok || item == itemPick && c.f.bag.Count(itemPick) > 0 || !c.f.bag.TakeAll(r.needs) {
		return false
	}
	c.f.bag.Add(item, 1)
	if item == itemPick {
		c.f.stock(itemPick)
	}
	return true
}

// itemOf is the inventory's item id, carried or not.
func (f *farm) itemOf(id string) backpack.Item {
	it, _ := carried{f}.Find(id)
	return it
}

// craftTabs are the names of the inventory's tabs.
func craftTabs() []string { return []string{i18n.T("Backpack"), i18n.T("Craft")} }
