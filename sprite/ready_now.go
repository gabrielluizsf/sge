package sprite

// Ready makes the image now if it is made on first use and has not been:
// a game that knows what it will show soon draws it a little at a time
// ahead, rather than all in the frame that first shows it.
func (s *Sprite) Ready() { s.ready() }

// IsReady reports whether the image is there to draw.
func (s *Sprite) IsReady() bool { return s.Image != nil || s.Make == nil }
