package farm

// tendedUnseen reports whether p is a garden of a town asleep, whose
// people go on watering it while no one sees them.
func (f *farm) tendedUnseen(p *plot) bool {
	return p.town != nil && !f.keeper.Awake(p.town)
}
