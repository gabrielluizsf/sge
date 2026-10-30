package farm

import "github.com/gabrielluizsf/sge/sprite"

// Body is how a person looks as they move: a clip of each of their
// animations, in each of the ways they may face on screen.
type Body interface {
	// Facings is how many ways the person may face on screen, 4 or 8, as
	// Facings names them.
	Facings() int
	// Clip is the animation anim, one of the Anim names, facing dir.
	Clip(anim, dir string) Clip
}

// Clip is an animation's frames, each a sprite whose pivot stands on the
// person's feet, and how many of them are shown a second. A clip has as
// many frames as its art has; AnimAct and AnimChop are shown once and
// should last ActSeconds, and the others loop.
type Clip struct {
	Frames []*sprite.Sprite
	FPS    float64
}

// The animations a person's Body has: standing, walking, using what they
// hold, chopping with the axe, and swimming.
const (
	AnimStand = "stand"
	AnimWalk  = "walk"
	AnimAct   = "act"
	AnimChop  = "chop"
	AnimSwim  = "swim"
)

// ActSeconds is how long a use of what one holds takes, chopping
// included, and so how long its clip should last.
const ActSeconds = actFor
