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

// The animations a person's Body has: standing, walking, running, using what they
// hold, chopping with the axe, forging at the anvil, hoeing, watering, picking a crop, swimming, diving into deep water from
// standing, rising out of it to stand, and sitting on a chair. The last two are shown once,
// and should last DiveSeconds and RiseSeconds.
const (
	AnimStand = "stand"
	AnimWalk  = "walk"
	AnimRun   = "run"
	AnimAct   = "act"
	AnimChop  = "chop"
	AnimForge = "forge"
	AnimHoe   = "hoe"
	AnimWater = "water"
	AnimPick  = "pick"
	AnimMine  = "mine"
	AnimEat   = "eat"
	AnimDrink = "drink"
	// Riding: getting on, sitting on the horse, and getting down.
	AnimMount    = "mount"
	AnimRide     = "ride"
	AnimDismount = "dismount"
	AnimWell     = "well"
	AnimDip      = "dip"
	AnimSwim     = "swim"
	AnimDive     = "dive"
	AnimRise     = "rise"
	AnimSit      = "sit"
	// Milking: sitting down to it, milking in a loop, and getting up.
	AnimMilkDown = "milkdown"
	AnimMilk     = "milk"
	AnimMilkUp   = "milkup"
)

// DiveSeconds and RiseSeconds are how long going into deep water and
// coming up out of it take, and so how long their clips should last.
const (
	DiveSeconds = 0.8
	RiseSeconds = 1.0
)

// ForgeSeconds is how long a smith's stroke at the anvil takes, and so
// how long its clip should last: the hammer raised, struck and let rest.
const ForgeSeconds = 1.0

// HoeSeconds, WaterSeconds and PickSeconds are the same for a stroke of
// the hoe, a pour from the can and the picking of a crop.
const (
	HoeSeconds   = 1.4
	WaterSeconds = 2.0
	PickSeconds  = 1.6
	MineSeconds  = 0.6
	EatSeconds   = 2.2
	DrinkSeconds = 2.6
	// MountSeconds is how long getting on a horse takes; DismountSeconds how
	// long getting down, sitting sideways and jumping off, done when asked for.
	MountSeconds    = 1.8
	DismountSeconds = 0.7
	WellSeconds     = 2.6
	DipSeconds      = 2.4
	// MilkDownSeconds is how long sitting down to milk takes, and getting
	// up after, and MilkSeconds the milking between them.
	MilkDownSeconds = 0.8
	MilkSeconds     = 4.2
)

// forgeImpact is how far into a stroke the hammer lands on the iron, as
// human.Forge strikes it.
const forgeImpact = 0.7

// ActSeconds is how long a use of what one holds takes, chopping
// included, and so how long its clip should last.
const ActSeconds = actFor
