package farm

import "github.com/gabrielluizsf/sge/pathfind"

// farmWork is a farmer's work: their garden first, picking, watering and
// sowing it from the doorstep; then selling what they have grown at the
// stall, or buying seeds there when they have none.
func (m *mind) farmWork(f *farm) *job {
	v := m.v
	k := v.kit
	if v.home == nil {
		return nil
	}
	tend := func(slot int) *job {
		return &job{to: pathfind.Cell(v.home.Door), face: "right", slot: slot, uses: 1, rest: 1}
	}
	if p := f.plotAt(v.home.Garden.Col, v.home.Garden.Row); p != nil {
		switch {
		case p.crop >= 0 && p.stage == ripe:
			return tend(-1)
		case p.thirsty():
			return tend(slotCan)
		case p.crop < 0:
			for c := len(crops) - 1; c >= 0; c-- {
				if k.bag.Count(seedsOf(c)) > 0 {
					return tend(c)
				}
			}
		}
	}
	seeds := 0
	for c := range crops {
		seeds += k.bag.Count(seedsOf(c))
	}
	switch {
	case worth(k) >= 20 || seeds == 0 && worth(k) > 0:
		return m.atStall(func(f *farm, at pathfind.Cell) { f.sell(k, f.center(at)) })
	case seeds == 0 && k.money >= crops[0].seedPrice:
		return m.atStall(func(f *farm, at pathfind.Cell) {
			c := v.luck.Intn(len(crops))
			for range 2 {
				if k.money >= crops[c].seedPrice {
					f.buy(k, c, f.center(at))
				}
			}
		})
	}
	return nil
}

// atStall is a trip to the town's stall, dealing there.
func (m *mind) atStall(deal func(f *farm, stall pathfind.Cell)) *job {
	s := pathfind.Cell(m.v.town.Stall)
	return &job{to: pathfind.Cell{Col: s.Col, Row: s.Row + 1}, face: "up", slot: -1, rest: 1.5, then: func(f *farm) { deal(f, s) }}
}
