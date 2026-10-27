package game

import "github.com/gabrielluizsf/antui/canvas"

// A Spreader is a script whose object is drawn as pieces, each at its own
// depth among the other objects of its Z: an isometric layer of trees,
// each in front of what stands behind it and behind what stands in front.
// Spread hands add each piece, the screen y its depth is and the painter
// of it, and is false when the object is drawn whole, as any other.
type Spreader interface {
	Spread(o *Object, add func(depth float64, paint func(*canvas.Canvas))) bool
}

// item is a piece of the world to draw: an object, or a piece of one,
// among the others of its Z by depth.
type item struct {
	z     int
	depth float64
	paint func(*canvas.Canvas)
}
