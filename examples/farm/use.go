package main

// use is k using what they hold on cell col, row: a ripe crop is picked
// whatever they hold; otherwise seeds are sown in empty soil, the can
// waters a thirsty crop, and the axe strikes a tree, or the anvil, as a
// hammer. In a room there is nothing to use. It reports whether
// it did anything.
func (f *farm) use(k *kit, col, row int) bool {
	if f.room != nil {
		return false
	}
	if k.slot == slotAxe && f.thingAt(col, row) == tileID["anvil"] {
		at := f.ground.CellCenter(col, row)
		f.sound("clang", &at)
		return true
	}
	p := f.plotAt(col, row)
	if p == nil || k.slot == slotAxe && !(p.crop >= 0 && p.stage == ripe) {
		return k.slot == slotAxe && f.chop(k, col, row)
	}
	at := f.ground.CellCenter(col, row)
	switch {
	case p.crop >= 0 && p.stage == ripe:
		c := p.crop
		k.bag.Add(cropOf(c), 1)
		p.crop, p.stage, p.watered, p.growing = -1, 0, false, 0
		f.float(at, "+1 "+crops[c].label, 'y')
		f.sound("harvest", &at)
	case p.crop < 0 && k.slot < slotCan && k.bag.Take(seedsOf(k.slot), 1):
		p.crop, p.stage = k.slot, 0
		f.sound("plant", &at)
	case k.slot == slotCan && p.thirsty():
		p.watered = true
		f.sound("pour", &at)
	case p.crop < 0 && k.slot < slotCan:
		f.float(at, "SEM SEMENTES", 'r')
		f.sound("deny", nil)
		return false
	default:
		return false
	}
	f.show(p)
	return true
}
