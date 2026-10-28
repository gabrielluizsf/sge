package sprite

// ready makes the image with Make, once, if there is none yet.
func (s *Sprite) ready() {
	if s.Image == nil && s.Make != nil {
		s.Image, s.Make = s.Make(), nil
	}
}
