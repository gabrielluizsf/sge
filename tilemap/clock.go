package tilemap

import "github.com/gabrielluizsf/sge/game"

// clock is the script a map runs on: it moves its animated tiles on,
// repainting only the cells whose frame changed, and only while the map
// can be seen.
type clock struct{ m *Map }

func (c clock) Update(o *game.Object, dt float64) {
	m := c.m
	if m.Tiles.animated == 0 {
		return
	}
	m.time += dt
	if !m.inView() {
		return // the cells catch up on the step they come into view
	}
	anims := m.Tiles.anims
	for i, id := range m.cells {
		if id != Empty && anims[id] != nil {
			m.paint(i%m.Cols, i/m.Cols)
		}
	}
}
