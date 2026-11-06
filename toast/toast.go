// Package toast is the timing and layout of small notices that slide in at
// the top of the screen, stay for a time that follows their length, and fade
// out, stacking downward. It draws nothing itself but BlitAlpha: the game
// asks where each toast is and how opaque, and paints it with its own art.
package toast

const (
	// EnterFor is how long a toast takes to slide in while fading in.
	EnterFor = 0.25
	// ExitFor is how long a toast takes to fade out while rising.
	ExitFor = 0.2
	// MaxShown is how many toasts are shown at once, older ones leave early.
	MaxShown = 3
	// Smooth is how fast, per second, a toast closes the gap to its place.
	Smooth = 18.0
)

// Toast is one notice: its text, how long it has been up, how long it stays
// before it starts to leave, and the row it is drawn on, which eases
// toward Target so a stack closes its gaps smoothly.
type Toast struct {
	Text   string
	Age    float64
	Life   float64
	Row    float64
	Target float64
}
