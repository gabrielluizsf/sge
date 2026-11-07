package farm

import "github.com/gabrielluizsf/sge/pathfind"

// presence is what a villager is doing at an hour of a day, and where:
// the chore their plan has them at then, and the job of it, the cell they
// stand on, the way they face and what they strike. It is worked out from
// the plan and the world alone, nothing simulated, so it costs the same
// for any hour of any day.
type presence struct {
	chore chore
	job   *job
}

// presenceAt is what v does at the hour of the day, or the zero presence
// when they have no chore then: they are asleep, at table, about town, or
// away on a visit, which whereabouts tells.
func (f *farm) presenceAt(v *villager, day int, hour float64) presence {
	plan := f.planOn(v, day)
	i := plan.index(hour)
	if i < 0 || plan[i].doing == choreVisit {
		return presence{}
	}
	c := plan[i].doing
	return presence{chore: c, job: f.siteOf(v, c, spellKey(v, day, i))}
}

// spellKey tells apart the spells of a day, for the choice of a place.
func spellKey(v *villager, day, spell int) uint32 {
	return mix(v.id, day) + uint32(spell)*977
}

// siteOf is the job of the chore c as it is done at the place key picks,
// among the places there are for it about the villager's town, or nil when
// there is none.
func (f *farm) siteOf(v *villager, c chore, key uint32) *job {
	switch c {
	case choreChop:
		return oneOf(f.sitesOf(f.chopQuarry(v), false), key)
	case choreMine:
		return oneOf(f.sitesOf(f.mineQuarry(v), false), key)
	case chorePlant:
		if v.home != nil {
			return &job{to: pathfind.Cell(v.home.Door), face: "right", slot: slotHoe, rest: 1}
		}
	case choreForge:
		if t := v.town; v.smith() && t != nil {
			return &job{to: pathfind.Cell{Col: t.Anvil.Col, Row: t.Anvil.Row + 1}, face: "up", slot: slotAxe, uses: forgeStrokes, rest: 1}
		}
	case choreTend:
		if v.town != nil {
			return dealerPost(f, v.town, key)
		}
	case choreTrade:
		if s := v.place(); s != nil {
			return &job{to: pathfind.Cell{Col: s.Stall.Col, Row: s.Stall.Row + 1}, face: "up", slot: -1, rest: 1.5}
		}
	}
	return nil
}

// oneOf is one of the jobs, by key, or nil when there is none.
func oneOf(jobs []*job, key uint32) *job {
	if len(jobs) == 0 {
		return nil
	}
	return jobs[int(key%uint32(len(jobs)))]
}
