package farm

import "github.com/gabrielluizsf/sge/progress"

// The XP of each deed: little for the common ones, more for a harvest
// that took days, a felled tree or a sale, and the most for an order or
// the errand.
const (
	xpSow      = 2
	xpMilk     = 4
	xpWolf     = 25
	xpWater    = 2
	xpHoe      = 1
	xpFell     = 12
	xpMine     = 10
	xpOrder    = 40
	xpErrand   = 100
	xpGift     = 5
	xpMarry    = 150
	xpChild    = 100
	xpLevelCut = 25 // money a level up pays, per level reached
)

// xpPick is the XP of picking crop c: more the longer it took to grow.
func xpPick(c int) int { return 4 + 2*crops[c].daysToRipe }

// xpSale is the XP of selling goods worth earned: a little for each $10,
// and at least 1.
func xpSale(earned int) int { return max(1, earned/10) }

// xpNeed is the XP a level takes.
func xpNeed(level int) int { return progress.Need(level) }
