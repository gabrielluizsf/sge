package farm

import "github.com/gabrielluizsf/sge/i18n"

// The night's hours: bed may be gone to from bedtime, and one still up at
// faintAt, a minute past midnight, faints where they stand, and all their
// money is lost. One wakes at wakeAt.
const (
	bedtime = 18.0
	faintAt = 1.0 / 60
	wakeAt  = 6.0
)

// sleep puts the player to bed in their house, if it is late enough: the
// night passes, the watered crops grow a day, and they wake
// in the morning.
func (f *farm) sleep() bool {
	if !f.clock.Between(bedtime, wakeAt) {
		f.float(f.player.Transform.Position, i18n.T("Too early"), 'r')
		f.sound("deny", nil)
		return false
	}
	f.rest()
	f.nextMorning()
	f.notify(i18n.T("Day %d: good morning!", f.clock.Day+1))
	return true
}

// nextMorning passes the night to wakeAt: the next day when it is still
// evening.
func (f *farm) nextMorning() {
	if f.clock.Hour >= wakeAt {
		f.clock.Day++
	}
	f.clock.Hour = wakeAt
	f.days()
}

// curfew makes the player faint if they are still up past midnight: they
// lose all their money, and wake at home in the morning.
func (f *farm) curfew() {
	if !f.clock.Between(faintAt, wakeAt) || f.fainted == f.clock.Day || f.question != nil {
		return
	}
	f.fainted = f.clock.Day
	lost := f.money
	f.money = 0
	f.nextMorning()
	f.fainted = f.clock.Day
	f.goHome()
	f.notify(i18n.T("You fainted! You lost $%d", lost))
	f.sound("deny", nil)
}

// goHome takes the player into their house, out of any other and of any
// conversation or shop.
func (f *farm) goHome() {
	f.talk, f.shop.open = nil, false
	if f.room != nil && f.room.own {
		return
	}
	if f.room != nil {
		f.leave()
	}
	f.enter(houseDoor{at: homeDoor, faceOut: "down", kind: 0, seed: seed, own: true})
}
