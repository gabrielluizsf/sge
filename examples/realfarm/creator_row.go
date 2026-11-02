package main

import "github.com/gabrielluizsf/sge/human"

// rowKind is how a row is changed.
type rowKind int

const (
	choiceRow rowKind = iota // arrows step through options
	sliderRow                // a value along a track
	buttonRow                // pressed
	textRow                  // typed into
)

// row is one line of the creation screen: its label, how it is changed
// and what it shows.
type row struct {
	label string
	kind  rowKind
	show  func(*charState) bool // nil: always shown

	text   func(*charState) string       // the value, for a choice
	swatch func(*charState) human.Colour // a colour beside the value, if set
	step   func(s *charState, d int)     // a choice stepped d (+1 or -1)
	slider *slider                       // a slider's range and place
	press  func(c *creator)              // a button pressed
}

// slider is a value from lo to hi that the row's field holds.
type slider struct {
	lo, hi float64
	field  func(*charState) *float64
	after  func(*charState) // done after a change, to fit the rest
}

// frac is the slider's place along the track, 0 to 1.
func (sl *slider) frac(s *charState) float64 {
	return (*sl.field(s) - sl.lo) / (sl.hi - sl.lo)
}

// set puts the slider at f, 0 to 1 along the track.
func (sl *slider) set(s *charState, f float64) {
	*sl.field(s) = sl.lo + min(max(f, 0), 1)*(sl.hi-sl.lo)
	if sl.after != nil {
		sl.after(s)
	}
}

// nudge moves the slider d (+1 or -1) of twenty steps along its track.
func (sl *slider) nudge(s *charState, d int) { sl.set(s, sl.frac(s)+float64(d)/20) }

// pressable reports whether pressing the row does something: a button
// is pressed, a text field takes the keyboard.
func (r *row) pressable() bool { return r.kind == buttonRow || r.kind == textRow }

// visible reports whether the row is shown for s.
func (r *row) visible(s *charState) bool { return r.show == nil || r.show(s) }
