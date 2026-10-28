package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// Picture is every tile of the map painted in its cell, as the cells show
// them now. It is made the first time it is wanted.
func (m *Map) Picture() *canvas.Canvas {
	if m.picture.Image == nil {
		img, err := canvas.NewCanvas(m.Cols*m.Tiles.Width, m.Rows*m.Tiles.Height)
		if err != nil {
			panic(err) // New checked the size
		}
		m.picture.Image = img
		for i, id := range m.painted {
			if id != Empty {
				copyTile(img, i%m.Cols*m.Tiles.Width, i/m.Cols*m.Tiles.Height, m.Tiles.Tiles[id].Image)
			}
		}
	}
	return m.picture.Image
}
