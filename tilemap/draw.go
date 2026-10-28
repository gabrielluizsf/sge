package tilemap

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// drawing is the script that draws a map, at the pixels drawing its whole
// picture as a sprite would put them: a scaled or turned map as its
// picture; one wholly opaque, copied from its picture in one go; one at
// least half see-through, blended from it in one go; any other tile by
// tile, copying opaque tiles and blending the rest.
type drawing struct{ m *Map }

func (drawing) Update(*game.Object, float64) {}

func (d drawing) Paint(o *game.Object, cv *canvas.Canvas) {
	m := d.m
	if g := o.Game(); g != nil && g.Projection != nil && !o.Screen {
		m.paintIso(cv)
		return
	}
	t := o.Transform
	if !o.Screen {
		t.Position = o.Game().ToScreen(t.Position)
	}
	if t.Scale != (vector.Vector3{X: 1, Y: 1, Z: 1}) || t.Rotation != vector.Identity() {
		m.Picture()
		m.picture.Draw(cv, t)
		return
	}
	// As sprite.Draw places a picture only moved.
	x0, y0 := -int(math.Floor(0.5-t.Position.X)), -int(math.Floor(0.5-t.Position.Y))
	if m.opaque() {
		cv.Copy(x0, y0, m.Picture())
		return
	}
	if 2*m.glassy >= len(m.cells) {
		cv.Blit(x0, y0, m.Picture())
		return
	}
	w, h := m.Tiles.Width, m.Tiles.Height
	for i, id := range m.painted {
		switch {
		case id == Empty:
		case m.Tiles.opaque[id]:
			cv.Copy(x0+i%m.Cols*w, y0+i/m.Cols*h, m.Tiles.Tiles[id].Image)
		default:
			cv.Blit(x0+i%m.Cols*w, y0+i/m.Cols*h, m.Tiles.Tiles[id].Image)
		}
	}
}
