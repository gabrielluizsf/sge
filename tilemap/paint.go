package tilemap

// paint draws into the map's picture, if it has one yet, the tile cell
// col, row shows now, unless it shows it already.
func (m *Map) paint(col, row int) {
	i := row*m.Cols + col
	id := m.cells[i]
	if id != Empty {
		id = m.Tiles.shown(id, col, row, m.time)
	}
	if id == m.painted[i] {
		return
	}
	m.count(m.painted[i], -1)
	m.count(id, 1)
	m.painted[i] = id
	img := m.picture.Image
	if img == nil {
		return // painted when the picture is first wanted
	}
	w, h := m.Tiles.Width, m.Tiles.Height
	if id == Empty {
		clearCell(img, col*w, row*h, w, h)
		return
	}
	copyTile(img, col*w, row*h, m.Tiles.Tiles[id].Image)
}
