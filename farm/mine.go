package farm

import "github.com/gabrielluizsf/sge/i18n"

// rockHits is how many strokes of the pickaxe break a rock, and fewest and
// most stones a broken rock gives.
const (
	rockHits   = 4
	stoneLeast = 2
	stoneMost  = 4
)

// swings reports whether k holds a tool swung at what stands on a cell,
// the axe or the pickaxe.
func (k *kit) swings() bool { return k.slot == slotAxe || k.slot == slotPick }

// mine is k's stroke of their pickaxe at the rock on cell col, row: after
// rockHits strokes it breaks, leaving nothing, giving k some stone, and
// grows back in time. It reports whether there was a rock to strike.
func (f *farm) mine(k *kit, col, row int) bool {
	rock := f.thingAt(col, row)
	if !isRock(rock) || k.bag.Count(itemPick) == 0 {
		return false
	}
	at := f.ground.CellCenter(col, row)
	key := cellKey(col, row)
	f.hits[key]++
	f.sound("clang", &at)
	if f.hits[key] < rockHits {
		return true
	}
	delete(f.hits, key)
	n := f.luck.Between(stoneLeast, stoneMost)
	k.bag.Add(stone, n)
	f.setThing(col, row, leftBy(rock))
	f.stumps = append(f.stumps, &stump{col: col, row: row, tree: rock, left: rockRegrowIn})
	f.earnBy(k, xpMine)
	f.float(at, i18n.T("+%d stone", n), 'y')
	if k == f.kit {
		f.notify(i18n.T("+%d stone", n))
	}
	f.sound("fell", &at)
	return true
}
