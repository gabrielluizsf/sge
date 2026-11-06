package farm

import "github.com/gabrielluizsf/sge/game"

// Update runs the world while the game plays, wherever the player is:
// stumps grow back; the villagers think of what to ask for; each new day
// is lived, the crops growing with it; and the player faints if they are
// up too late.
func (f *farm) Update(o *game.Object, dt float64) {
	f.drift(dt)
	f.vitalsTick(dt)
	f.ring(dt)
	f.fade(dt)
	f.xpTick(dt)
	f.regrow(dt)
	for _, v := range f.people {
		v.ask(dt)
	}
	f.days()
	f.curfew()
}
