package farm

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/sge/backpack"
)

// The items that are not crops, seeds or wood: the tools in hand.
const (
	itemCan = "can"
	itemHoe = "hoe"
	itemAxe = "axe"
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
	return items
}

// Find is the item id, whether the player has it or not.
func (c carried) Find(id string) (backpack.Item, bool) {
	switch {
	case id == itemCan || id == itemHoe || id == itemAxe:
		return c.tool(id), true
	case id == wood:
		return c.wood(), true
	case strings.HasSuffix(id, seedsSuffix):
		i, ok := cropIndex(strings.TrimSuffix(id, seedsSuffix))
		return c.seeds(i), ok
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
		return backpack.Item{ID: id, Name: "Regador", Icon: "icon_can", Count: f.inCan(), Hand: true,
			Info: []string{"Rega a terra para as plantas crescerem", fmt.Sprint("Água: ", f.inCan(), "/", canCapacity), "Encha na água"}}
	case itemHoe:
		return backpack.Item{ID: id, Name: "Enxada", Icon: "icon_hoe", Count: -1, Hand: true,
			Info: []string{"Limpa a planta seca do canteiro"}}
	}
	a := axes[max(f.axe, 0)]
	return backpack.Item{ID: id, Name: "Machado de " + a.label, Icon: "icon_axe_" + a.name, Count: -1, Hand: true,
		Info: []string{"Corta árvores", fmt.Sprint(a.hits, " golpes por árvore, até ", a.most, " de madeira")}}
}

// seeds is the packet of crop i's seeds.
func (c carried) seeds(i int) backpack.Item {
	cr := crops[i]
	return backpack.Item{ID: seedsOf(i), Name: "Sementes de " + cr.label, Icon: "icon_" + cr.seedIcon, Count: c.f.bag.Count(seedsOf(i)), Hand: true,
		Info: []string{fmt.Sprint("Planta no canteiro; amadurece em ", cr.daysToRipe, " dias"), fmt.Sprint("Compra por $", cr.seedPrice)}}
}

// harvest is crop i, picked.
func (c carried) harvest(i int) backpack.Item {
	cr := crops[i]
	return backpack.Item{ID: cropOf(i), Name: cr.label, Icon: "icon_" + cr.cropIcon, Count: c.f.bag.Count(cropOf(i)),
		Info: []string{fmt.Sprint("Vende por $", cr.sellPrice), "Presente para quem gosta"}}
}

// wood is the wood from the trees.
func (c carried) wood() backpack.Item {
	return backpack.Item{ID: wood, Name: "Madeira", Icon: "icon_" + wood, Count: c.f.bag.Count(wood),
		Info: []string{fmt.Sprint("Vende por $", woodPrice), "Presente para o ferreiro e os lenhadores"}}
}
