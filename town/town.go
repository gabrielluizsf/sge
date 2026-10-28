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
	// Seed is the town's own, for what is up to the game: its people.
	Seed uint32
}

// House is a house of a town: its Size by Size cells from Col, Row, its
// kind of look, its doorstep in the yard in front of it, and the yard's
// other cell, a garden.
type House struct {
	Col, Row int
	Kind     int
	Door     Cell
	Garden   Cell
}

// The plan's measures, in cells: a house's side; how far the main street
// runs each way from the plaza's middle; the plaza's half width and half
// height; and how far from the middle, along the street, each house of a
// side starts.
const (
	Size       = 2
	Reach      = 9
	plazaW     = 2
	plazaH     = 1
	houseDepth = Size + 1 // a house and its yard, beside the street
)

// slots are where houses may stand along a side of the street: their
// first column, from the plaza's middle.
var slots = []int{-8, -5, 3, 6}
