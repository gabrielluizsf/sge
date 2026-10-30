package main

import "fmt"

// chop is k's stroke of their axe at the tree on cell col, row: after as
// many strokes as their axe takes, it falls, leaving a stump and giving k
// some wood. It reports whether there was a tree to strike.
func (f *farm) chop(k *kit, col, row int) bool {
	if f.thingAt(col, row) != tileID["tree0"] {
		return false
	}
	at := f.ground.CellCenter(col, row)
	if k.axe == noAxe {
		f.float(at, "Sem machado", 'r')
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
	f.stumps = append(f.stumps, &stump{col: col, row: row, left: regrowIn})
	f.float(at, fmt.Sprint("+", n, " madeira"), 'y')
	f.sound("fell", &at)
	return true
}
