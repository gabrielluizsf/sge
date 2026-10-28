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
		mx, my := Reach+2, houseDepth+2
		src := random.New(terrain.Hash(l.Seed+1, r.Col, r.Row))
		center = Cell{r.Col*l.Region + mx + src.Intn(l.Region-2*mx), r.Row*l.Region + my + src.Intn(l.Region-2*my)}
	}
	t := &Town{Region: r, Center: center, Stall: Cell{center.Col, center.Row - plazaH - 1}, Anvil: Cell{center.Col - plazaW, center.Row - plazaH - 1}, Forge: Cell{center.Col - plazaW - 1, center.Row - plazaH - 1}, Seed: terrain.Hash(l.Seed+2, r.Col, r.Row)}
	src := random.New(t.Seed)
	for _, side := range []int{-1, 1} {
		for _, dx := range slots {
			if src.Float() < 0.2 {
				continue // an empty lot
			}
			h := House{Col: center.Col + dx, Kind: src.Intn(max(l.Kinds, 1))}
			yard := center.Row + side
			if side < 0 {
				h.Row = yard - Size
			} else {
				h.Row = yard + 1
			}
			h.Door, h.Garden = Cell{h.Col, yard}, Cell{h.Col + 1, yard}
			t.Houses = append(t.Houses, h)
		}
	}
	return t
}
