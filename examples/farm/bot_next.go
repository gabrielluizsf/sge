package main

import "github.com/gabrielluizsf/sge/game"

// next is the bot's next task: pick what is ripe, water what is thirsty,
// sell a harvest, sow seeds it has, and buy seeds when it has none; or
// nothing, while the crops grow.
func (b *bot) next(f *farm) *task {
	plots := f.plots[:min(b.tends, len(f.plots))]
	tend := func(p *plot, slot int) *task {
		return &task{to: standIn(f, p.col, p.row-1), face: "down", do: func(f *farm, o *game.Object) bool {
			f.slot = slot
			return true
		}}
	}
	for _, p := range plots {
		if p.crop >= 0 && p.stage == ripe {
			return tend(p, f.slot)
		}
	}
	for _, p := range plots {
		if p.thirsty() {
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
		for c := len(crops) - 1; c >= 0; c-- {
			if f.bag.Count(seedsOf(c)) > 0 {
				return tend(p, c)
			}
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
				c := 1
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
