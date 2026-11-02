package farm

import "github.com/gabrielluizsf/sge/pathfind"

// farmWork is a farmer's work: their garden first, picking, clearing
// what withered, watering and sowing it from the doorstep, and going to
// water to fill the can when it is empty; then selling what they have
// grown at the stall, or buying seeds there when they have none.
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
		case p.ripened():
			return tend(-1)
		case p.dead:
			return tend(slotHoe)
		case p.thirsty() && k.inCan() == 0:
			return m.fetchWater(f)
		case p.thirsty():
			return tend(slotCan)
		case p.crop < 0:
			if c, ok := k.sowing(); ok {
				j := tend(slotSeeds)
				j.seed = c
				return j
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

// fetchWater is a trip to the nearest water to fill the can; with none
// in reach, the can is filled at the well of home, the doorstep.
func (m *mind) fetchWater(f *farm) *job {
	v := m.v
	if !m.looked {
		m.shore, m.found = f.springNear(pathfind.Cell(v.home.Door))
		m.looked = true
	}
	if !m.found {
		v.kit.fill()
		return nil
	}
	return &job{to: m.shore.stand, face: m.shore.face, slot: slotCan, uses: 1, rest: 1}
}

// atStall is a trip to the town's stall, dealing there.
func (m *mind) atStall(deal func(f *farm, stall pathfind.Cell)) *job {
	s := pathfind.Cell(m.v.town.Stall)
	return &job{to: pathfind.Cell{Col: s.Col, Row: s.Row + 1}, face: "up", slot: -1, rest: 1.5, then: func(f *farm) { deal(f, s) }}
}
