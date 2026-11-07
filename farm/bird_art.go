package farm

import "github.com/gabrielluizsf/sge/sprite"

// BirdArt is what an Art may also have to show birds: the farm has none
// when its look cannot draw them.
type BirdArt interface {
	// Bird is the sprite of a bird as look says, its pivot at its middle.
	Bird(look BirdLook) *sprite.Sprite
}

// BirdLook is how a bird is seen: flying or perched, the eighth of a turn
// it faces on the screen, 0 toward the right and on clockwise, and the
// frame of its wings, which a perched bird has none of.
type BirdLook struct {
	Flying  bool
	Heading int
	Frame   int
}

// BirdFrames is how many frames a flying bird's wings go through; a
// BirdArt wraps any frame number round it.
const BirdFrames = 8
