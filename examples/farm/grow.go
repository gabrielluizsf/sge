package farm

import "github.com/gabrielluizsf/sge/game"

// Update runs the world while the game plays, wherever the player is:
// every watered crop grows, after its crop's stage time a stage further
// on and dry again until it is watered once more; stumps grow back; the
// villagers think of what to ask for; and the player faints if they are
// up too late.
func (f *farm) Update(o *game.Object, dt float64) {
	for _, p := range f.plots {
		if p.crop < 0 || !p.watered || p.stage >= ripe {
			continue
		}
		if p.growing += dt; p.growing >= crops[p.crop].stageTime {
			p.stage++
			p.watered, p.growing = false, 0
			f.show(p)
		}
	}
	f.drift(dt)
	f.fade(dt)
	f.regrow(dt)
	for _, v := range f.people {
		v.ask(dt)
	}
	f.curfew()
}
