package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
)

// woolDays is how many days a sheep's wool takes to grow back.
const woolDays = 3

// sheepBy is the player's sheep nearest o, within reach, out of doors.
func (f *farm) sheepBy(o *game.Object) *creature {
	if f.herd == nil || f.room != nil || f.milking != nil || f.hands.swims {
		return nil
	}
	var near *creature
	best := milkReach
	for _, c := range f.herd.crits {
		if d := c.pos.Sub(feet(o)).Length(); c.kind == "sheep" && c.pen.town == nil && d < best {
			near, best = c, d
		}
	}
	return near
}

// woolIn is how many days more sheep c's wool takes to grow, 0 when it may
// be shorn; milked is the day she last gave it, plus one.
func (f *farm) woolIn(c *creature) int {
	if c.milked == 0 {
		return 0
	}
	return max(0, c.milked-1+woolDays-f.clock.Day)
}

// sheepNear is the interaction of a sheep within reach: shear her, or be
// told how long her wool takes to grow.
func (f *farm) sheepNear(o *game.Object) (words string, do func(), ok bool) {
	c := f.sheepBy(o)
	switch {
	case c == nil:
		return "", nil, false
	case f.woolIn(c) > 0:
		return i18n.T("Wool grows back in %d days", f.woolIn(c)), func() { f.sound("deny", nil) }, true
	}
	return i18n.T("E: shear the sheep"), func() {
		if f.shear(c) {
			clip, seconds := useClip(actPick)
			f.hands.acting = seconds
			f.player.Animator.Play(f.hands.pose(f.player) + clip)
		}
	}, true
}

// shear has the player take sheep c's wool, if it has grown, and reports
// whether they did.
func (f *farm) shear(c *creature) bool {
	if f.woolIn(c) > 0 {
		return false
	}
	c.milked = f.clock.Day + 1
	f.bag.Add(wool, 1)
	at := feet(f.player)
	at.Y -= tile
	f.float(at, "+1 "+i18n.T("Wool"), 'y')
	f.sound("harvest", &at)
	f.earnBy(f.kit, xpMilk)
	return true
}
