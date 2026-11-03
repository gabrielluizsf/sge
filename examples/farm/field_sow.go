package farm

// fieldSow is a villager on duty taking free soil of the player's farm:
// a bare plot nobody keeps, to sow their seeds in, if they have seeds
// and the caps leave room. Of the free plots one is picked at random.
func (m *mind) fieldSow(f *farm) *job {
	v := m.v
	crop, ok := v.kit.sowing()
	if !ok || !f.onDuty(v) || !f.roomInField(v) {
		return nil
	}
	var free []*job
	for _, p := range f.farmPlots() {
		if p.crop >= 0 || p.owner != nil {
			continue
		}
		if at, face, ok := f.standBy(p.col, p.row); ok {
			free = append(free, &job{to: at, face: face, slot: slotSeeds, seed: crop, uses: 1, rest: 1})
		}
	}
	if len(free) == 0 {
		return nil
	}
	return free[v.luck.Intn(len(free))]
}
