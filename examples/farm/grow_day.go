package farm

// growCrops is the morning's work on the soil: a crop watered yesterday
// grows a day further, and goes dry again; one still growing that went
// the whole day unwatered withers. A ripe crop waits to be picked
// without water, and a withered one waits for the hoe. A garden of a town
// asleep, far from the player, is watered by its people unseen.
func (f *farm) growCrops() {
	for _, p := range f.plots {
		if p.crop < 0 || p.dead || p.stage >= ripe {
			continue
		}
		if p.watered || f.tendedUnseen(p) {
			p.days++
			p.stage = stageOf(crops[p.crop], p.days)
			p.watered = false
		} else {
			p.dead = true
		}
		f.show(p)
	}
}

// pick takes the ripe crop from p: it is gone, or, if it bears again,
// back to the stage it regrows from.
func (f *farm) pick(p *plot) {
	c := crops[p.crop]
	if c.regrowDays == 0 {
		p.clear()
		return
	}
	p.days = c.daysToRipe - c.regrowDays
	p.stage, p.watered = stageOf(c, p.days), false
}
