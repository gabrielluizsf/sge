package tilemap

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// Spread hands a Sorted map's tiles in view to be drawn each where its
// cell stands among the objects of its Z, in a projected world.
func (d drawing) Spread(o *game.Object, add func(float64, func(*canvas.Canvas))) bool {
	m := d.m
	if !m.Sorted || o.Game() == nil || o.Game().Projection == nil {
		return false
	}
	c0, r0, c1, r1 := m.isoCells()
	for r := r0; r <= r1; r++ {
		for c := c0; c <= c1; c++ {
			if m.painted[r*m.Cols+c] != Empty {
				add(m.isoDepth(c, r), func(cv *canvas.Canvas) { m.paintCell(cv, c, r) })
			}
		}
	}
	return true
}
