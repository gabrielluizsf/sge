package tilemap

// turns reports whether tile id shows different frames as time goes: it is
// animated, and not by where its cell lies, as a Parity animation is.
func (s *TileSet) turns(id int) bool {
	return id >= 0 && s.anims != nil && s.anims[id] != nil && !s.anims[id].Parity
}
