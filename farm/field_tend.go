package farm

// fieldTend is a villager's care of their crops on the player's farm,
// when they have any: picking what is ripe, clearing what withered and
// watering what is dry, fetching water when the can is empty.
func (m *mind) fieldTend(f *farm) *job {
	v := m.v
	for _, p := range f.farmPlots() {
		if p.owner != v {
			continue
		}
		slot := -1
		switch {
		case p.ripened():
		case p.dead:
			slot = slotHoe
		case p.thirsty() && v.kit.inCan() == 0:
			return m.fetchWater(f)
		case p.thirsty():
			slot = slotCan
		default:
			continue
		}
		if at, face, ok := f.standBy(p.col, p.row); ok {
			return &job{to: at, face: face, slot: slot, uses: 1, rest: 1}
		}
	}
	return nil
}
