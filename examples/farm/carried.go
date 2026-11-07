package farm

import (
	"strings"

	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// The items that are not crops, seeds or wood: the tools in hand.
const (
	itemCan  = "can"
	itemHoe  = "hoe"
	itemAxe  = "axe"
	itemPick = "pickaxe"
)

// seedsSuffix ends the name of an item of seeds.
const seedsSuffix = "_seeds"

// carried is what the player carries, as the inventory screen lists it:
// the tools, the seeds, the harvest and the wood.
type carried struct{ f *farm }

// Items are the tools, then each packet of seeds, each crop picked and the
// wood, of those the player has.
func (c carried) Items() []backpack.Item {
	f := c.f
	items := []backpack.Item{c.tool(itemCan), c.tool(itemHoe)}
	if f.axe != noAxe {
		items = append(items, c.tool(itemAxe))
	}
	if f.bag.Count(itemPick) > 0 {
		items = append(items, c.tool(itemPick))
	}
	for i := range crops {
		if f.bag.Count(seedsOf(i)) > 0 {
			items = append(items, c.seeds(i))
		}
	}
	for i := range crops {
		if f.bag.Count(cropOf(i)) > 0 {
			items = append(items, c.harvest(i))
		}
	}
	if f.bag.Count(wood) > 0 {
		items = append(items, c.wood())
	}
	if f.bag.Count(stone) > 0 {
		items = append(items, c.stone())
	}
	for _, g := range goods {
		if f.bag.Count(g.item) > 0 {
			items = append(items, c.good(g))
		}
	}
	for _, name := range furnishings {
		if f.bag.Count(name) > 0 {
			items = append(items, c.furniture(name))
		}
	}
	return items
}

// Find is the item id, whether the player has it or not.
func (c carried) Find(id string) (backpack.Item, bool) {
	switch {
	case id == itemCan || id == itemHoe || id == itemAxe || id == itemPick:
		return c.tool(id), true
	case id == wood:
		return c.wood(), true
	case id == stone:
		return c.stone(), true
	case isFurnishing(id):
		return c.furniture(id), true
	case strings.HasSuffix(id, seedsSuffix):
		i, ok := cropIndex(strings.TrimSuffix(id, seedsSuffix))
		return c.seeds(i), ok
	}
	if g, ok := goodOf(id); ok {
		return c.good(g), true
	}
	i, ok := cropIndex(id)
	return c.harvest(i), ok
}

// cropIndex is the index of the crop called name.
func cropIndex(name string) (int, bool) {
	for i := range crops {
		if crops[i].name == name {
			return i, true
		}
	}
	return 0, false
}

// seedsCrop is the crop whose seeds the item id is.
func seedsCrop(id string) (int, bool) {
	if !strings.HasSuffix(id, seedsSuffix) {
		return 0, false
	}
	return cropIndex(strings.TrimSuffix(id, seedsSuffix))
}

// tool is the tool id, with what it does.
func (c carried) tool(id string) backpack.Item {
	f := c.f
	switch id {
	case itemCan:
		return backpack.Item{ID: id, Name: i18n.T("Watering can"), Icon: "icon_can", Count: f.inCan(), Hand: true,
			Info: []string{i18n.T("Waters the soil for plants to grow"), i18n.T("Water: %d/%d", f.inCan(), canCapacity), i18n.T("Fill at the water")}}
	case itemHoe:
		return backpack.Item{ID: id, Name: i18n.T("Hoe"), Icon: "icon_hoe", Count: -1, Hand: true,
			Info: []string{i18n.T("Clears the withered plant from the plot")}}
	case itemPick:
		return backpack.Item{ID: id, Name: i18n.T("Pickaxe"), Icon: "icon_pickaxe", Count: -1, Hand: true,
			Info: []string{i18n.T("Breaks rocks into stone"), i18n.T("%d hits per rock", rockHits)}}
	}
	a := axes[max(f.axe, 0)]
	return backpack.Item{ID: id, Name: i18n.T("%s axe", i18n.T(a.label)), Icon: "icon_axe_" + a.name, Count: -1, Hand: true,
		Info: []string{i18n.T("Chops trees"), i18n.T("%d hits per tree, up to %d wood", a.hits, a.most)}}
}

// seeds is the packet of crop i's seeds.
func (c carried) seeds(i int) backpack.Item {
	cr := crops[i]
	return backpack.Item{ID: seedsOf(i), Name: i18n.T("%s seeds", i18n.T(cr.label)), Icon: "icon_" + cr.seedIcon, Count: c.f.bag.Count(seedsOf(i)), Hand: true,
		Info: []string{i18n.T("Plant in the plot; ripens in %d days", cr.daysToRipe), i18n.T("Buy for $%d", cr.seedPrice)}}
}

// harvest is crop i, picked.
func (c carried) harvest(i int) backpack.Item {
	cr := crops[i]
	return backpack.Item{ID: cropOf(i), Name: i18n.T(cr.label), Icon: "icon_" + cr.cropIcon, Count: c.f.bag.Count(cropOf(i)), Hand: true,
		Info: []string{i18n.T("Sells for $%d", cr.sellPrice), i18n.T("Eat: +%d food", int(cr.nutrition)), i18n.T("Gift for those who like it")}}
}

// wood is the wood from the trees.
func (c carried) wood() backpack.Item {
	return backpack.Item{ID: wood, Name: i18n.T("Wood"), Icon: "icon_" + wood, Count: c.f.bag.Count(wood),
		Info: []string{i18n.T("Sells for $%d", woodPrice), i18n.T("Gift for the smith and the lumberjacks")}}
}
