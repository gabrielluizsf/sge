package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// The opacity of a tile, worked out the first time it is asked for.
const (
	opacityUnknown int8 = iota
	seeThrough
	solidColour
)

// isOpaque reports whether tile id has no transparent pixel; Empty has
// nothing but. A tile not drawn yet is drawn to find out.
func (s *TileSet) isOpaque(id int) bool {
	if id < 0 || id >= len(s.opaque) {
		return false
	}
	if s.opaque[id] == opacityUnknown {
		s.opaque[id] = seeThrough
		if isOpaque(s.picture(id)) {
			s.opaque[id] = solidColour
		}
	}
	return s.opaque[id] == solidColour
}

// picture is tile id's image, drawn now if it is made when first wanted.
func (s *TileSet) picture(id int) *canvas.Canvas { return s.Tiles[id].Picture() }
