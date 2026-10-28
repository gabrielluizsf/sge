package tilemap

// isOpaque reports whether tile id has no transparent pixel; Empty has
// nothing but.
func (s *TileSet) isOpaque(id int) bool {
	return id >= 0 && id < len(s.opaque) && s.opaque[id]
}
