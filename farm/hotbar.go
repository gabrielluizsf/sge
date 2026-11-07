package farm

import (
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/i18n"
)

// hotbarSlots is how many things the player may have at hand, each on a
// number key; inventoryCols how many items a row of the inventory has.
const (
	hotbarSlots   = 6
	inventoryCols = 8
)

// newPack is the player's inventory screen and hotbar over what they
// carry, with the first seeds, the can and the hoe at hand.
func newPack(f *farm) *backpack.Pack {
	p := backpack.New(carried{f}, hotbarSlots, inventoryCols)
	p.Scale = scale
	p.Title = i18n.T("Backpack")
	p.Hint = i18n.T("Enter: pick up   1-6: shortcut   Drag: mouse   I/Tab/Esc: close")
	p.Hotbar = defaultHotbar()
	p.Rows, p.Tabs, p.Crafts = 4, craftTabs(), craftBook{f}
	return p
}

// defaultHotbar is what a new farmer has at hand: the seeds of the first
// crop on 1, the can on 2 and the hoe on 4, the axe's slot waiting for it.
func defaultHotbar() []string {
	h := make([]string, hotbarSlots)
	h[0], h[1], h[3] = seedsOf(0), itemCan, itemHoe
	return h
}

// hold takes in hand what is on hotbar slot pos, if anything: the seeds
// of that crop, the can, the hoe or the axe.
func (f *farm) hold(pos int) bool {
	if pos < 0 || pos >= len(f.pack.Hotbar) {
		return false
	}
	id := f.pack.Hotbar[pos]
	switch c, seeds := seedsCrop(id); {
	case seeds:
		f.slot, f.seed = slotSeeds, c
	case id == itemCan:
		f.slot = slotCan
	case id == itemHoe:
		f.slot = slotHoe
	case id == itemAxe:
		f.slot = slotAxe
	case id == itemPick && f.bag.Count(itemPick) > 0:
		f.slot = slotPick
	default:
		if g, isGood := goodOf(id); isGood {
			if g.nutrition <= 0 {
				return false
			}
			f.slot, f.food = slotFood, id
			break
		}
		c, crop := cropIndex(id)
		if !crop {
			return false
		}
		f.slot, f.seed, f.food = slotFood, c, ""
	}
	f.pack.Active = pos
	return true
}

// keepHand takes in hand again what the hotbar slot in hand holds, after
// the player rearranged the hotbar.
func (f *farm) keepHand() { f.hold(f.pack.Active) }

// stock puts item id on the first empty slot of the hotbar if it is on
// none, as a tool just got is.
func (f *farm) stock(id string) {
	if f.pack.Slot(id) >= 0 {
		return
	}
	for i, h := range f.pack.Hotbar {
		if h == "" {
			f.pack.Bind(id, i)
			return
		}
	}
}
