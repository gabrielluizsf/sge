// Package town lays out towns over an endless world: the world is cut
// into square regions, each with a town or none as a seed decides, and
// each town a main street with a plaza in the middle, a market stall on
// it, and houses of several kinds along the street, each with a yard in
// front. The same seed lays out the same towns, wherever one looks first.
package town

// Cell is a cell of the world.
type Cell struct{ Col, Row int }

// Town is a town's plan.
type Town struct {
	// Region is the region it is the town of.
	Region Cell
	// Center is the middle of its plaza, on its main street.
	Center Cell
	// Houses are its houses, west to east, north side first.
	Houses []House
	// Stall is the market stall's cell, at the north edge of the plaza:
	// one buys standing on the cell south of it.
	Stall Cell
	// Anvil is the smith's anvil's cell, west of the stall at the plaza's
	// edge: the smith works on the cell south of it.
	Anvil Cell
	// Forge is the smith's furnace, west of the anvil: it is fed from the
	// cell south of it.
	Forge Cell
	// Well is the corner, the north west one, of the cells south of the
	// plaza's middle where a well may stand: the game says whether it has
	// one, and how large.
	Well Cell
	// Reach is how far from the plaza's middle the main street runs each
	// way, and Depth how many cells the houses and their yards reach from
	// the street's middle on each side, as the houses' sizes make them.
	Reach, Depth int
	// Seed is the town's own, for what is up to the game: its people.
	Seed uint32
	// road is the road to the town it leads to, once linked.
	road   []Rect
	to     *Town
	linked bool
}

// House is a house of a town: its Width by Depth cells from Col, Row, its
// kind of look, its doorstep in the yard in front of it, in the middle
// of its front, and the yard's next cell, a garden.
type House struct {
	Col, Row     int
	Width, Depth int
	Kind         int
	Door         Cell
	Garden       Cell
}

// The plan's measures, in cells: the plaza's half width and half height;
// how far from the plaza's edge the first house of a side of the street
// stands; and the footprint a house has when the layout says none.
const (
	plazaW      = 2
	plazaH      = 1
	defaultLead = 4
	defaultW    = 2
)
