package tilemap

// SetSolid makes the tiles ids solid, or not.
func (s *TileSet) SetSolid(solid bool, ids ...int) {
	for _, id := range ids {
		if id >= 0 && id < len(s.solid) {
			s.solid[id] = solid
		}
	}
}
