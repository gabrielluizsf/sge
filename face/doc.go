// Package face keeps a talking face alive: its eyes blink at a person's
// pace, every two to six seconds, a blink taking about a sixth of a
// second, and its mouth takes the shape of the letters being said, each
// held long enough to be seen, shut between words and when all is said.
// It says how the face looks; an art draws it.
package face

// Look is how a face is at a moment: its eyelids, LidsOpen to LidsShut,
// and its mouth's shape.
type Look struct {
	Lids  int
	Mouth Shape
}

// How far shut the eyelids are.
const (
	LidsOpen = iota
	LidsHalf
	LidsShut
)

// Shape is the mouth's shape as a letter is said.
type Shape int

// The mouth's shapes: shut, for m, b and p and between words; ajar, for
// most consonants; open, for a, e and i; and round, for o and u.
const (
	Shut Shape = iota
	Ajar
	Open
	Round
)
