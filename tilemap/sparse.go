package tilemap

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// sparse draws a map with no Sprite, one under half full, tile by tile:
// only the cells holding one, at the pixels drawing its whole picture
// would put them. A scaled or turned map is drawn as its picture.
type sparse struct{ m *Map }

func (sparse) Update(*game.Object, float64) {}

func (s sparse) Paint(o *game.Object, cv *canvas.Canvas) {
	m := s.m
	if m.Sprite != nil || m.filled == 0 {
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
	w, h := m.Tiles.Width, m.Tiles.Height
	for i, id := range m.painted {
		if id != Empty {
			cv.Blit(x0+i%m.Cols*w, y0+i/m.Cols*h, m.Tiles.Tiles[id].Image)
		}
	}
}
