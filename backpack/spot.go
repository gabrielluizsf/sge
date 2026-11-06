package backpack

// Area is a part of the screen with slots: the grid of what is carried,
// the hotbar, or the chest beside the grid.
type Area int

// The areas.
const (
	Nowhere Area = iota
	Grid
	Bar
	Chest
)

// Spot is a slot: its area and its place in it.
type Spot struct {
	Area  Area
	Index int
}

// Event is what a key or the mouse did, for the game to answer with a
// sound.
type Event int

// The events.
const (
	Nothing Event = iota
	Opened
	Closed
	Moved   // the cursor went to another slot
	Picked  // an item was taken up
	Bound   // an item was put on a hotbar slot
	Unbound // an item was taken off the hotbar
	Dropped // an item was let go where it came from
	Refused // an item can not go there
	Crafted // a recipe was made
	Stored  // an item went between the grid and the chest
	Tabbed  // another tab was shown
)
