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
	// Footprint is how many cells across, and deep, a house of a kind
	// stands on; with none, two by two. The towns are spaced by the
	// largest, so no two houses touch, and as many houses line a side of
	// the street, up to two each way, as the region has room for.
	Footprint func(kind int) Cell
	// Gap is the cells left between two houses along the street, at least
	// one.
	Gap int
	// Lead is how many cells of plaza the first houses keep from its
	// middle, at least four.
	Lead int
	// Margin is how many cells more the houses keep from the plaza's
	// props, the stall, the forge and the well, whose pictures stand
	// taller than their cells: it widens the Lead.
	Margin int
	Fixed  map[Cell]Cell
	Skip   func(rx, ry int) bool
	// Roads are roads laid by hand, besides those between the towns; Keep
	// are the areas, the farm, the roads between towns keep off.
	Roads, Keep []Rect
	towns       map[int]*Town // by region key; nil where there is none
	// roads are the rectangles of the roads near a region, by its key; last
	// is the answer for the region asked about last.
	roads          map[int][]Rect
	last           []Rect
	lastRx, lastRy int
	lastOK         bool
	// biggest is the largest footprint, once lotKnown.
	biggest  Cell
	lotKnown bool
}
