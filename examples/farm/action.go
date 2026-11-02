package farm

// action is what a use did: the thing the person is then shown doing.
type action int

// The actions a use may be: nothing, picking a ripe crop, watering a plot,
// filling the can, hoeing, felling at a tree, striking at the anvil, and
// sowing.
const (
	actNone action = iota
	actPick
	actWater
	actFill
	actHoe
	actChop
	actForge
	actSow
)
