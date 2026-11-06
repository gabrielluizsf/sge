package tilemap

import "github.com/gabrielluizsf/sge/sprite"

// PicturesOf are the sprites a cell holding tile id shows over time: its
// own and, if it is animated, its frames'.
func (s *TileSet) PicturesOf(id int) []*sprite.Sprite {
	if id < 0 || id >= len(s.Tiles) {
		return nil
	}
	out := []*sprite.Sprite{s.Tiles[id]}
	if id < len(s.anims) && s.anims[id] != nil {
		for _, f := range s.anims[id].Frames {
			if f >= 0 && f < len(s.Tiles) && f != id {
				out = append(out, s.Tiles[f])
			}
		}
	}
	return out
}
