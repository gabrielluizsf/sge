package farm

// use is k using what they hold on cell col, row: a ripe crop is picked
// whatever they hold; otherwise seeds are sown in empty soil, the can
// waters a thirsty crop or is filled at water, the hoe clears a withered
// plant, and the axe strikes a tree, or the anvil, as a hammer. In a room
// there is nothing to use. It reports whether it did anything.
func (f *farm) use(k *kit, col, row int) bool { return f.useAs(k, col, row) != actNone }

// useAs is use, reporting what it did.
func (f *farm) useAs(k *kit, col, row int) action {
	if f.room != nil {
		return actNone
	}
	if k.slot == slotAxe && f.thingAt(col, row) == tileID["anvil"] {
		f.soundIn(cue{name: "clang", at: f.ground.CellCenter(col, row), in: forgeImpact * ForgeSeconds})
		return actForge
	}
	if k.slot == slotCan && f.isWater(col, row) {
		if f.refill(k, col, row) {
			return actFill
		}
		return actNone
	}
	p := f.plotAt(col, row)
	if p == nil || k.slot == slotAxe && !p.ripened() {
		if k.slot == slotAxe && f.chop(k, col, row) {
			return actChop
		}
		return actNone
	}
	at := f.ground.CellCenter(col, row)
	did := actNone
	switch {
	case p.ripened():
		did = actPick
		c := p.crop
		k.bag.Add(cropOf(c), 1)
		f.pick(p)
		f.float(at, "+1 "+crops[c].label, 'y')
		f.sound("harvest", &at)
	case p.dead && k.slot == slotHoe:
		did = actHoe
		p.clear()
		f.sound("plant", &at)
	case p.crop < 0 && k.slot == slotSeeds && f.sow(k, p):
		did = actSow
		f.sound("plant", &at)
	case k.slot == slotCan && p.thirsty():
		if !k.pour() {
			if k == f.kit {
				f.notify("O regador está vazio — encha na água")
			}
			f.sound("deny", nil)
			return actNone
		}
		did = actWater
		p.watered = true
		f.sound("pour", &at)
	case p.crop < 0 && k.slot == slotSeeds:
		f.float(at, "Sem sementes", 'r')
		f.sound("deny", nil)
		return actNone
	default:
		return actNone
	}
	f.show(p)
	return did
}

// sow sows k's seeds in the bare plot p, if they have any.
func (f *farm) sow(k *kit, p *plot) bool {
	c, ok := k.sowing()
	if !ok || !k.bag.Take(seedsOf(c), 1) {
		return false
	}
	p.crop, p.stage, p.days = c, 0, 0
	return true
}

// refill fills k's can from the water on cell col, row, and reports
// whether there was room in it.
func (f *farm) refill(k *kit, col, row int) bool {
	if k.used == 0 {
		return false
	}
	k.fill()
	at := f.ground.CellCenter(col, row)
	f.sound("pour", &at)
	if k == f.kit {
		f.float(at, "Regador cheio", 'g')
	}
	return true
}
