package farm

// useClip is the clip a use is shown by, as the suffix of its name, and
// how long it takes: by what the use did, not by what is in hand, so a
// ripe crop picked with the can in hand is a pick, not a pour. A stroke
// of the axe, of the hammer at the anvil, of the hoe, a pour, a pick, or
// the act of anything else.
func useClip(a action) (string, float64) {
	switch a {
	case actForge:
		return "_forge", ForgeSeconds
	case actChop:
		return "_chop", actFor
	case actMine:
		return "_mine", MineSeconds
	case actHoe:
		return "_hoe", HoeSeconds
	case actWater:
		return "_water", WaterSeconds
	case actPick:
		return "_pick", PickSeconds
	}
	return "_act", actFor
}
