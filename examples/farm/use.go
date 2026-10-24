package main

// use is the farmer using what they hold on the plot at col, row: a ripe
// crop is picked whatever they hold; otherwise seeds are sown in empty
// soil and the can waters a thirsty crop. It reports whether it did
// anything.
func (f *farm) use(col, row int) bool {
	p := f.plotAt(col, row)
	if p == nil {
		return false
	}
	at := f.ground.CellCenter(col, row)
	switch {
	case p.crop >= 0 && p.stage == ripe:
		c := p.crop
		f.harvest[c]++
		p.crop, p.stage, p.watered, p.growing = -1, 0, false, 0
		f.float(at, "+1 "+crops[c].label, 'y')
	case p.crop < 0 && f.slot < slotCan && f.seeds[f.slot] > 0:
		f.seeds[f.slot]--
		p.crop, p.stage = f.slot, 0
	case f.slot == slotCan && p.thirsty():
		p.watered = true
	case p.crop < 0 && f.slot < slotCan:
		f.float(at, "SEM SEMENTES", 'r')
		return false
	default:
		return false
	}
	f.show(p)
	return true
}
