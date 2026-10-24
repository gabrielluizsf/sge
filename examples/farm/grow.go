package main

import "github.com/gabrielluizsf/sge/game"

// Update grows every watered crop: after its crop's stage time it is a
// stage further on, and dry again until it is watered once more.
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
}
