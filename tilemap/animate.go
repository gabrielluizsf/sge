package tilemap

// Animate makes tile id show the frames of a: cells holding id show them
// in turn as the world steps.
func (s *TileSet) Animate(id int, a Animation) {
	if id < 0 || id >= len(s.Tiles) || len(a.Frames) == 0 {
		return
	}
	if s.anims == nil {
		s.anims = make([]*Animation, len(s.Tiles))
	}
	if s.anims[id] == nil {
		s.animated++
	}
	s.anims[id] = &a
}
