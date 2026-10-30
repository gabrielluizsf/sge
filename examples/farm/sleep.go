package farm

import "fmt"

// The night's hours: bed may be gone to from bedtime, and one still up at
// faintAt, a minute past midnight, faints where they stand, and all their
// money is lost. One wakes at wakeAt.
const (
	bedtime = 18.0
	faintAt = 1.0 / 60
	wakeAt  = 6.0
)

// sleep puts the player to bed in their house, if it is late enough: the
// night passes, the watered crops grow a stage overnight, and they wake
// in the morning.
func (f *farm) sleep() bool {
	if !f.clock.Between(bedtime, wakeAt) {
		f.float(f.player.Transform.Position, "Ainda é cedo", 'r')
		f.sound("deny", nil)
		return false
	}
	f.nextMorning()
	f.notify(fmt.Sprint("Dia ", f.clock.Day+1, ": bom dia!"))
	return true
}

// nextMorning passes the night to wakeAt: the next day when it is still
// evening.
func (f *farm) nextMorning() {
	if f.clock.Hour >= wakeAt {
		f.clock.Day++
	}
	f.clock.Hour = wakeAt
	for _, p := range f.plots {
		if p.watered && p.stage < ripe {
			p.stage++
			p.watered, p.growing = false, 0
			f.show(p)
		}
	}
}

// curfew makes the player faint if they are still up past midnight: they
// lose all their money, and wake at home in the morning.
func (f *farm) curfew() {
	if !f.clock.Between(faintAt, wakeAt) || f.fainted == f.clock.Day {
		return
	}
	f.fainted = f.clock.Day
	lost := f.money
	f.money = 0
	f.talk, f.shop.open = nil, false
	f.nextMorning()
	f.fainted = f.clock.Day
	if f.room == nil || !f.room.own {
		if f.room != nil {
			f.leave()
		}
		f.enter(houseDoor{at: homeDoor, faceOut: "down", kind: 0, seed: seed, own: true})
	}
	f.notify(fmt.Sprint("Você desmaiou! Perdeu $", lost))
	f.sound("deny", nil)
}
