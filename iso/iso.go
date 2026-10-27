// Package iso sees a game's world from a corner, as isometric games do:
// the world stays a flat grid, where things move and collide as in any
// game, and only its picture is turned, each square cell drawn as a
// diamond twice as wide as it is tall. A Projection is a game.Projector;
// a game given one draws its tile maps as diamonds and its objects in
// depth order, nearer over farther.
package iso

// Projection draws a world cell Cell pixels square as a diamond Width
// pixels wide and Height tall: the cell's corner at the world's origin is
// the diamond's top, its x axis runs down and to the right, its y axis
// down and to the left.
type Projection struct {
	Cell, Width, Height float64
}

// New is the projection of cells cell pixels square onto diamonds width
// by height.
func New(cell, width, height float64) Projection {
	return Projection{Cell: cell, Width: width, Height: height}
}
