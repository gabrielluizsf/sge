package tilemap

// Cell is the size of a tile's cell in the world: CellWidth by
// CellHeight, or a tile's picture's where they are not set.
func (t *TileSet) Cell() (width, height int) {
	if t.CellWidth > 0 && t.CellHeight > 0 {
		return t.CellWidth, t.CellHeight
	}
	return t.Width, t.Height
}
