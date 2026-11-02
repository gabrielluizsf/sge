package town

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/terrain"
)

// Of is the town of region rx, ry, or nil where there is none. A town is
// planned the first time it is asked for and the same plan given after.
func (l *Layout) Of(rx, ry int) *Town {
	r := Cell{rx, ry}
	k, ok := key(r)
	if !ok {
		return nil
	}
	if t, ok := l.towns[k]; ok {
		return t
	}
	if l.towns == nil {
		l.towns = map[int]*Town{}
	}
	t := l.plan(r)
	l.towns[k] = t
	return t
}

// plan is region r's town, or nil.
func (l *Layout) plan(r Cell) *Town {
	center, fixed := l.Fixed[r]
	if !fixed {
		if l.Skip != nil && l.Skip(r.Col, r.Row) || terrain.Chance(l.Seed, r.Col, r.Row) >= l.Chance {
			return nil
		}
		// The plaza's middle anywhere the whole town stays in the region.
		mx, my := l.Reach()+2, l.Depth()+2
		src := random.New(terrain.Hash(l.Seed+1, r.Col, r.Row))
		center = Cell{r.Col*l.Region + mx + src.Intn(l.Region-2*mx), r.Row*l.Region + my + src.Intn(l.Region-2*my)}
	}
	t := &Town{Region: r, Center: center, Reach: l.Reach(), Depth: l.Depth(), Well: Cell{center.Col - 1, center.Row + plazaH + 1}, Stall: Cell{center.Col, center.Row - plazaH - 1}, Anvil: Cell{center.Col - plazaW, center.Row - plazaH - 1}, Forge: Cell{center.Col - plazaW - 1, center.Row - plazaH - 1}, Seed: terrain.Hash(l.Seed+2, r.Col, r.Row)}
	src := random.New(t.Seed)
	for _, side := range []int{-1, 1} {
		for _, s := range l.slots() {
			if src.Float() < 0.2 {
				continue // an empty lot
			}
			s.side = side
			t.Houses = append(t.Houses, l.house(t, s, src.Intn(max(l.Kinds, 1))))
		}
	}
	return t
}

// house is the house of the kind in slot s, its yard between it and the
// street.
func (l *Layout) house(t *Town, s slot, kind int) House {
	fp := l.size(kind)
	h := House{Col: t.Center.Col + s.start, Width: fp.Col, Depth: fp.Row, Kind: kind}
	if s.west { // the house's far side stands at the lot's, whatever its width
		h.Col += l.lot().Col - fp.Col
	}
	yard := t.Center.Row + s.side
	h.Row = yard + 1
	if s.side < 0 {
		h.Row = yard - fp.Row
	}
	door := h.Col + (fp.Col-1)/2
	h.Door, h.Garden = Cell{door, yard}, Cell{door + 1, yard}
	return h
}
