package tilemap

import "github.com/gabrielluizsf/antui/canvas"

// trimmed is a tile's picture cut to its pixels that show, and where in
// the whole picture that cut starts: an isometric tile is mostly the air
// above its diamond, and blending air is time lost.
type trimmed struct {
	img    *canvas.Canvas
	dx, dy int
	spans  [][2]int // its rows' opaque runs, when it has nothing else
}

// trim is tile id cut to what shows, cut the first time it is wanted.
func (t *TileSet) trim(id int) trimmed {
	if len(t.trims) < len(t.Tiles) {
		t.trims = append(t.trims, make([]trimmed, len(t.Tiles)-len(t.trims))...)
	}
	if t.trims[id].img == nil {
		t.trims[id] = cut(t.Tiles[id].Image)
	}
	return t.trims[id]
}

// cut is img's box of pixels that are not wholly transparent, copied.
func cut(img *canvas.Canvas) trimmed {
	x0, y0, x1, y1 := img.Width, img.Height, -1, -1
	for y := range img.Height {
		for x := range img.Width {
			if img.At(x, y).A() != 0 {
				x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x), max(y1, y)
			}
		}
	}
	if x1 < 0 {
		x0, y0, x1, y1 = 0, 0, 0, 0
	}
	out, err := canvas.NewCanvas(x1-x0+1, y1-y0+1)
	if err != nil {
		panic(err)
	}
	out.Copy(-x0, -y0, img)
	return trimmed{img: out, dx: x0, dy: y0, spans: spansOf(out)}
}
