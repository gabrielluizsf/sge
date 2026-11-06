package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// use is k using what they hold on cell col, row: a ripe crop is picked
// whatever they hold; otherwise seeds are sown in empty soil, the can
// waters a thirsty crop or is filled at water, the hoe clears a withered
// plant, the axe strikes a tree, or the anvil, as a hammer, and the pickaxe
// a rock. In a room
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
	if k.slot == slotFood && !f.ripeAt(col, row) {
		return f.eat(k)
	}
	if k.slot == slotCan && f.isWater(col, row) {
		if f.refill(k, col, row) {
			return actFill
		}
		return actNone
	}
	p := f.plotAt(col, row)
	if p == nil || k.swings() && !p.ripened() {
		switch {
		case k.slot == slotAxe && f.chop(k, col, row):
			return actChop
		case k.slot == slotPick && f.mine(k, col, row):
			return actMine
		}
		return actNone
	}
	if !f.mayTend(k, p) {
		f.refuse(k, p)
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
		f.earnBy(k, xpPick(c))
		f.float(at, "+1 "+i18n.T(crops[c].label), 'y')
		f.sound("harvest", &at)
	case p.dead && k.slot == slotHoe:
		did = actHoe
		p.clear()
		f.earnBy(k, xpHoe)
		f.sound("plant", &at)
	case p.crop < 0 && k.slot == slotSeeds && f.sow(k, p):
		did = actSow
		f.earnBy(k, xpSow)
		f.sound("plant", &at)
	case k.slot == slotCan && p.thirsty():
		if !k.pour() {
			if k == f.kit {
				f.notify(i18n.T("The watering can is empty — fill it at the water"))
			}
			f.sound("deny", nil)
			return actNone
		}
		did = actWater
		f.earnBy(k, xpWater)
		p.watered = true
		f.sound("pour", &at)
	case p.crop < 0 && k.slot == slotSeeds:
		f.float(at, i18n.T("No seeds"), 'r')
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
	f.claim(k, p)
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
		f.float(at, i18n.T("Watering can full"), 'g')
	}
	return true
}
