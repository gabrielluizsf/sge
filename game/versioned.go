package game

// A Versioned is a Painter or Spreader that can tell when what it draws
// has changed, so the game may leave a frame it would draw the same
// undrawn (see Frame). A script that is not a Versioned is taken to draw
// something new every frame.
type Versioned interface {
	// Version is a number that is the same for as long as the script draws
	// the same pixels, the camera and its object apart, which the game
	// looks at itself, and another whenever it would draw others.
	Version() uint64
}
