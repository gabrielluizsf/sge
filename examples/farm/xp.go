package farm

import "github.com/gabrielluizsf/sge/examples/i18n"

// glowFor is how long, in seconds, the bar and the ring glow after a level up.
const glowFor = 1.5

// earn gives the player n XP: a "+n XP" rising from them (but not over a
// conversation), and, when it
// takes them to a new level, news, a sound, a glow and a little money.
func (f *farm) earn(n int) {
	if n <= 0 {
		return
	}
	if f.talk == nil { // in a conversation the bar says it
		f.float(f.player.Transform.Position, i18n.T("+%d XP", n), 'g')
	}
	for range f.prog.Add(n) {
		f.money += xpLevelCut * f.prog.Level
		f.notify(i18n.T("Level %d!", f.prog.Level))
		f.sound("level_up", nil)
		f.glow = glowFor
	}
}

// earnBy is earn when k is the player's own kit: the villagers' work is
// not the player's deed.
func (f *farm) earnBy(k *kit, n int) {
	if k == f.kit {
		f.earn(n)
	}
}

// xpTick lets the bar catch up with the XP, smoothly, and the glow fade.
// The bar's value is levels reached plus the fraction of the next.
func (f *farm) xpTick(dt float64) {
	f.glow = max(0, f.glow-dt)
	f.shown += (f.xpValue() - f.shown) * min(1, 6*dt)
	if d := f.xpValue() - f.shown; d < 0.0005 && d > -0.0005 {
		f.shown = f.xpValue()
	}
}

// xpValue is how far along the player is: the level they have reached,
// from 0, plus the part of the next that they have.
func (f *farm) xpValue() float64 {
	return float64(f.prog.Level-1) + float64(f.prog.XP)/float64(xpNeed(f.prog.Level))
}
