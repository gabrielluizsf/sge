package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// paintIso draws a projected map's cells in view farthest first, a
// diagonal at a time, so a nearer tile is drawn over a farther one.
func (m *Map) paintIso(cv *canvas.Canvas) {
	c0, r0, c1, r1 := m.isoCells()
	for d := c0 + r0; d <= c1+r1; d++ {
		for c := max(c0, d-r1); c <= min(c1, d-r0); c++ {
			m.paintCell(cv, c, d-c)
		}
	}
}

// paintCell draws cell col, row's tile, if it has one and it shows on
// cv, where it goes in a projected world.
func (m *Map) paintCell(cv *canvas.Canvas, col, row int) {
	id := m.painted[row*m.Cols+col]
	if id == Empty {
		return
	}
	t := m.Tiles.trim(id)
	x, y := m.isoAt(col, row)
	x, y = x+t.dx, y+t.dy
	if x >= cv.Width || y >= cv.Height || x+t.img.Width <= 0 || y+t.img.Height <= 0 {
		return
	}
	if t.spans != nil {
		copySpans(cv, t, x, y)
		return
	}
	cv.Blit(x, y, t.img)
}
