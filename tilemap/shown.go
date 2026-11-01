package tilemap

import "math"

// shown is the tile cell col, row holding id shows t seconds in: id
// itself, or the frame of its animation.
func (s *TileSet) shown(id, col, row int, t float64) int {
	if s.anims == nil || id < 0 || s.anims[id] == nil {
		return id
	}
	a := s.anims[id]
	n := len(a.Frames)
	if a.Parity {
		return a.Frames[(col&1+(row&1)<<1)%n]
	}
	at := t * a.FPS
	if a.Staggered {
		at += float64(uint32(col*73856093^row*19349663)%997) / 997 * float64(n)
	}
	return a.Frames[int(math.Floor(at))%n]
}
