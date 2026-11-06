package farm

import (
	"fmt"
	"strings"
)

// knownItem reports whether item is one a bag or an order may hold: wood,
// a crop, or a packet of seeds.
func knownItem(item string) bool {
	if item == wood || item == stone || item == itemPick || isFurnishing(item) {
		return true
	}
	for c := range crops {
		if item == cropOf(c) || item == seedsOf(c) {
			return true
		}
	}
	return false
}

// hasStack reports whether l has a stack of item.
func hasStack(l []stackState, item string) bool {
	for _, s := range l {
		if s.Item == item {
			return true
		}
	}
	return false
}

// checkStacks is nil if l is a few stacks of known items, no two of the
// same, each a bounded number.
func checkStacks(what string, l []stackState) error {
	if len(l) > maxStacks {
		return fmt.Errorf("save: %s has %d stacks", what, len(l))
	}
	seen := map[string]bool{}
	for _, s := range l {
		if !knownItem(s.Item) || seen[s.Item] {
			return fmt.Errorf("save: %s has a bad item %q", what, s.Item)
		}
		seen[s.Item] = true
		if err := inRange(what+" "+s.Item, s.N, 1, maxCount); err != nil {
			return err
		}
	}
	return nil
}

// checkKit is nil if k is a kit the game could have.
func checkKit(what string, k kitState) error {
	return firstError(
		inRange(what+" money", k.Money, 0, maxMoney),
		checkStacks(what+" bag", k.Bag),
		inRange(what+" slot", k.Slot, slotSeeds, slotPick),
		inRange(what+" seed", k.Seed, 0, len(crops)-1),
		inRange(what+" axe", k.Axe, noAxe, len(axes)-1),
		inRange(what+" sold", k.Sold, 0, maxMoney),
		inRange(what+" used", k.Used, 0, maxCount),
	)
}

// checkOrder is nil if o is an order someone could ask.
func checkOrder(what string, o orderState) error {
	return firstError(checkStacks(what, o.Want), inRange(what+" reward", o.Reward, 0, maxMoney))
}

// checkHotbar is nil if h is a hotbar the game could have, with the slot in hand.
func checkHotbar(s *snapshot) error {
	h := s.Hotbar
	if len(h) != hotbarSlots {
		return fmt.Errorf("save: hotbar has %d slots, not %d", len(h), hotbarSlots)
	}
	seen := map[string]bool{}
	for _, id := range h {
		_, seeds := seedsCrop(id)
		_, crop := cropIndex(id)
		ok := id == "" || seeds || crop || isFurnishing(id) || id == itemCan || id == itemHoe || id == itemAxe && s.Player.Axe != noAxe || id == itemPick && hasStack(s.Player.Bag, itemPick)
		if !ok || strings.TrimSpace(id) != id || seen[id] && id != "" {
			return fmt.Errorf("save: bad hotbar item %q", id)
		}
		seen[id] = true
	}
	return inRange("active slot", s.Active, 0, hotbarSlots-1)
}
