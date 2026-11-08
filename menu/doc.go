// Package menu is the screens a game shows when it is paused: a List of
// lines to choose from with the arrows or the mouse, which scrolls when
// it has more than fits, as a pause menu or a list of credits is, and a
// Reader of a long text, as a licence is, scrolled by the arrows, the page
// keys or the wheel. Neither holds pictures: the game draws them through
// an Art.
package menu

import "github.com/gabrielluizsf/sge/backpack"

// Art is how the game draws the menus: panels, text and, if need be,
// pictures; the same the backpack uses.
type Art = backpack.Art

// Point and Rect are pixels, and parts of the screen.
type (
	Point = backpack.Point
	Rect  = backpack.Rect
)

// Input is what the player did this frame: arrows as steps, Enter or
// space, Escape, the page keys as steps of a page, the wheel's notches
// (positive scrolling down), and the mouse: where it is and whether its
// button went down.
type Input struct {
	Dy     int
	Page   int
	Wheel  int
	Select bool
	Cancel bool
	Mouse  Point
	Press  bool
}

// Result is what an input did to a screen.
type Result int

// The results.
const (
	Nothing   Result = iota
	Moved            // the choice or the scroll changed
	Chosen           // a line was chosen: see List.Choice
	Cancelled        // the player left the screen
)
