package tilemap

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// paintIso draws a projected map's cells that show, farthest first, a
// diagonal at a time, so a nearer tile is drawn over a farther one.
func (m *Map) paintIso(cv *canvas.Canvas) {
	r := m.isoCells()
	m.catchUp(r)
	for d := r.c0 + r.r0; d <= r.c1+r.r1; d++ {
		for c := max(r.c0, d-r.r1); c <= min(r.c1, d-r.r0); c++ {
			m.paintCell(cv, c, d-c)
		}
	}
}

// paintCell draws cell col, row's tile, if it has one and it shows on
// cv within its clip, where it goes in a projected world.
func (m *Map) paintCell(cv *canvas.Canvas, col, row int) {
	id := m.painted[row*m.Cols+col]
	if id == Empty {
		return
	}
	t := m.Tiles.trim(id)
	x, y := m.isoAt(col, row, t)
	x, y = x+t.dx, y+t.dy
	clip := cv.Clip
	if x >= clip.X+clip.Width || y >= clip.Y+clip.Height || x+t.img.Width <= clip.X || y+t.img.Height <= clip.Y {
		return
	}
	if t.h > m.Tiles.Height && game.Hides(game.Placed{Img: t.img, X: x, Y: y}, m.isoDepth(col, row), m.Game().Covers()) {
		game.BlendFaded(cv, game.Placed{Img: t.img, X: x, Y: y}) // a tall tile over what is kept in sight
		return
	}
	if t.spans != nil {
		copySpans(cv, t, place{x, y})
		return
	}
	cv.Blit(x, y, t.img)
}
