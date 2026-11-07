package farm

import "github.com/gabrielluizsf/sge/sprite"

// AnimalArt is what an Art may also have to show the farm's animals: the
// farm has none when its look cannot draw them.
type AnimalArt interface {
	// Animal is the sprite of an animal as look says, its pivot at its
	// feet, or nil if the look has no such picture.
	Animal(look AnimalLook) *sprite.Sprite
	// Pace is how fast kind walks, in cells a second, for its feet not to
	// slide, and how many frames of its walk the look shows a second.
	Pace(kind string) (speed, fps float64)
	// Milker is where one milking a cow sits, by her feet on the screen,
	// in pixels: how far behind her middle, and how far nearer the viewer,
	// with her head to the screen's right or left and they facing it too,
	// for the hands to come to her teats.
	Milker() (back, near float64)
	// Frames is how many pictures the kind has of a motion it plays once,
	// a wolf's "bite" or "fall".
	Frames(kind, motion string) int
}

// AnimalLook is how an animal is seen: its kind, "cow", "sheep", "horse"
// or "wolf"; its coat, 0 for its first; what it does, "stand", "walk" or
// "graze", or a wolf's "bite" and "fall" (a look without that picture shows
// it standing); the eighth of a
// turn it faces on the screen, 0 toward the right and on clockwise; and
// the frame of the motion, which a look wraps round its own.
type AnimalLook struct {
	Kind    string
	Coat    int
	Motion  string
	Heading int
	Frame   int
}
