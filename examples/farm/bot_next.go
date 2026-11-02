package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
)

// next is the bot's next task: pick what is ripe, clear what withered,
// water what is thirsty, going to the water to fill the can when it is
// empty, sell a harvest, sow seeds it has, and buy seeds when it has none;
// or nothing, while the crops grow.
func (b *bot) next(f *farm) *task {
	plots := f.plots[:min(b.tends, len(f.plots))]
	tend := func(p *plot, slot int) *task {
		return &task{to: standIn(f, p.col, p.row-1), face: "down", do: func(f *farm, o *game.Object) bool {
			f.slot = slot
			return true
		}}
	}
	for _, p := range plots {
		if p.ripened() {
			return tend(p, f.slot)
		}
	}
	for _, p := range plots {
		if p.dead {
			return tend(p, slotHoe)
		}
	}
	for _, p := range plots {
		switch {
		case p.thirsty() && f.inCan() == 0:
			return b.fetchWater(f)
		case p.thirsty():
			return tend(p, slotCan)
		}
	}
	if worth(f.kit) > 0 {
		return b.atStall(f, func(f *farm) { f.choose(f.kit, &shop{choice: sellChoice, at: farmStall}) })
	}
	for _, p := range plots {
		if p.crop >= 0 {
			continue
		}
		if c, ok := f.anySeeds(); ok {
			f.seed = c
			return tend(p, slotSeeds)
		}
	}
	empty := 0
	for _, p := range plots {
		if p.crop < 0 {
			empty++
		}
	}
	if empty > 0 && f.money >= crops[0].seedPrice {
		return b.atStall(f, func(f *farm) {
			// The dearer seeds while they can be had, for every empty plot.
			for range empty {
				c := botCrop
				if f.money < crops[c].seedPrice {
					c = 0
				}
				if f.money < crops[c].seedPrice {
					return
				}
				f.choose(f.kit, &shop{choice: c, at: farmStall})
			}
		})
	}
	return nil
}

// atStall is a task at the stall's counter, doing deal there.
func (b *bot) atStall(f *farm, deal func(f *farm)) *task {
	return &task{to: standIn(f, stallCol, stallRow+1), face: "up", do: func(f *farm, o *game.Object) bool {
		deal(f)
		return false
	}}
}

// botCrop is the crop the bot buys seeds of while it can pay for them,
// otherwise the first: the tomato, which bears again.
const botCrop = 3

// fetchWater is a task at the water's edge, filling the can; the bot
// looks for the water near the field once.
func (b *bot) fetchWater(f *farm) *task {
	if !b.looked {
		b.shore, b.found = f.springOnFarm(pathfind.Cell{Col: f.plots[0].col, Row: f.plots[0].row})
		b.looked = true
	}
	if !b.found {
		f.fill()
		return nil
	}
	return &task{to: standIn(f, b.shore.stand.Col, b.shore.stand.Row), face: b.shore.face, do: func(f *farm, o *game.Object) bool {
		f.slot = slotCan
		return true
	}}
}
