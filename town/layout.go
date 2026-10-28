package town

// Layout is where the towns are: a town in a region of Region by Region
// cells with the given Chance, from Seed. Fixed regions have a town for
// sure, its plaza's middle where they say; Skip, when set, keeps towns
// out of the regions it reports, where something else is.
type Layout struct {
	Seed   uint32
	Region int
	Chance float64
	// Kinds is how many kinds of houses there are to choose from.
	Kinds int
	Fixed map[Cell]Cell
	Skip  func(rx, ry int) bool
	towns map[int]*Town // by region key; nil where there is none
}
