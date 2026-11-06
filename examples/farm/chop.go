package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// chop is k's stroke of their axe at the tree on cell col, row: after as
// many strokes as their axe takes, it falls, leaving a stump and giving k
// some wood. It reports whether there was a tree to strike.
func (f *farm) chop(k *kit, col, row int) bool {
	tree := f.thingAt(col, row)
	if !isTree(tree) {
		return false
	}
	at := f.ground.CellCenter(col, row)
	if k.axe == noAxe {
		f.float(at, i18n.T("No axe"), 'r')
		f.sound("deny", nil)
		return false
	}
	key := cellKey(col, row)
	f.hits[key]++
	f.sound("chop", &at)
	a := axes[k.axe]
	if f.hits[key] < a.hits {
		return true
	}
	delete(f.hits, key)
	n := f.luck.Between(1, a.most)
	k.bag.Add(wood, n)
	f.setThing(col, row, tileID["stump"])
	f.stumps = append(f.stumps, &stump{col: col, row: row, tree: tree, left: regrowIn})
	f.earnBy(k, xpFell)
	f.float(at, i18n.T("+%d wood", n), 'y')
	f.sound("fell", &at)
	return true
}
