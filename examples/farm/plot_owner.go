package farm

// villagerOf is the villager whose kit k is, nil for none.
func (f *farm) villagerOf(k *kit) *villager {
	for _, v := range f.people {
		if v.kit == k {
			return v
		}
	}
	return nil
}

// mayTend reports whether k may work plot p: a town's garden as ever; on
// the farm, the player the soil nobody owns, and a villager their own
// crop, or bare soil nobody owns, to sow.
func (f *farm) mayTend(k *kit, p *plot) bool {
	switch {
	case p.town != nil:
		return true
	case k == f.kit:
		return p.owner == nil
	case p.owner != nil:
		return p.owner.kit == k
	}
	return p.crop < 0 && k.slot == slotSeeds
}

// refuse tells the player, when it was they who tried, whose crop p holds.
func (f *farm) refuse(k *kit, p *plot) {
	if k != f.kit || p.owner == nil {
		return
	}
	f.notify("Esta plantação é de " + p.owner.name)
	f.sound("deny", nil)
}

// claim makes the villager whose kit sowed p its owner, on the farm.
func (f *farm) claim(k *kit, p *plot) {
	if p.town == nil && k != f.kit {
		p.owner = f.villagerOf(k)
	}
}

// release gives back to the player what v kept on the farm.
func (f *farm) release(v *villager) {
	for _, p := range f.plots {
		if p.owner == v {
			p.owner = nil
		}
	}
}

// ownerLabel is " de <name>" for a crop that is a villager's, else nothing.
func ownerLabel(p *plot) string {
	if p.owner == nil {
		return ""
	}
	return " de " + p.owner.name
}
