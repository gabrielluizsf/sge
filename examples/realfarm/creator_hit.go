package main

// rect is a part of the screen, in pixels.
type rect struct{ X, Y, W, H int }

// has reports whether the point is in r.
func (r rect) has(x, y int) bool { return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H }

// hitKind is what a part of the screen does when pressed.
type hitKind int

const (
	hitCategory hitKind = iota
	hitRow
	hitLess
	hitMore
	hitTrack
	hitDropper
)

// hit is a part of the screen the mouse acts on, set when it is painted.
type hit struct {
	at       rect
	kind     hitKind
	cat, row int
}
