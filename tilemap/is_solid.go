package tilemap

// IsSolid reports whether the tile id is solid. Empty is not.
func (s *TileSet) IsSolid(id int) bool {
	return id >= 0 && id < len(s.solid) && s.solid[id]
}
