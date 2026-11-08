package tilemap

import "github.com/gabrielluizsf/sge/game"

// clock is the script a map runs on: it moves its animated tiles on,
// repainting only the cells whose frame changed, and only those the view
// shows. A cell that comes into view later is brought up to the map's
// time when it is first drawn (see catchUp), so nothing it shows is stale.
type clock struct{ m *Map }

func (c clock) Update(_ *game.Object, dt float64) {
	m := c.m
	if m.Tiles.animated == 0 {
		return
	}
	m.time += dt
	if m.Tiles.turning == 0 {
		return // what is animated by where a cell lies never turns over
	}
	m.fresh = cellRange{}
	if r := m.seen(); r.valid {
		m.refresh(r)
		m.fresh = r
	}
}
