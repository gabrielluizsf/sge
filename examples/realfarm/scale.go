package realfarm

import "github.com/gabrielluizsf/sge/iso"

// Cell is a cell's side in the world, and a metre: the models are made in
// metres, and the world is a grid of one-metre cells.
const Cell = 32

// Projection is the game's: each metre a diamond 64 by 32 pixels.
var Projection = iso.New(Cell, 2*Cell, Cell)
