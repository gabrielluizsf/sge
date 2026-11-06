package farm

// tendedUnseen reports whether p is a garden of a town whose people go on
// watering it while no one sees them: a town asleep, or one awake but out
// of sight, whose people stand still until they are seen.
func (f *farm) tendedUnseen(p *plot) bool {
	return p.town != nil && (!f.keeper.Awake(p.town) || !f.plotInSight(p))
}
