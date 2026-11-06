package vignette

import "github.com/gabrielluizsf/antui/canvas"

// Paint draws the vignette over cv, of its size, as strong as opacity, from
// 0 to 1, keeping cv's own alpha. A canvas of another size or whose pixels
// cannot be reached directly gets none.
func (v *Vignette) Paint(cv *canvas.Canvas, opacity float64) {
	if opacity <= 0 || cv.Pixels == nil || cv.Width != v.Width || cv.Height != v.Height {
		return
	}
	op := uint32(min(opacity, 1)*256 + 0.5)
	ink := uint32(v.cfg.Color)
	inkRB, inkG := ink&0xFF00FF, ink&0xFF00
	for y := range v.Height {
		n := v.reach[y]
		if n == 0 {
			continue
		}
		row := cv.Pixels[y*cv.Stride : y*cv.Stride+v.Width]
		mask := v.alpha[y*v.Width : (y+1)*v.Width]
		blendRow(row[:n], mask[:n], [3]uint32{op, inkRB, inkG})
		if rest := max(n, v.Width-n); rest < v.Width {
			blendRow(row[rest:], mask[rest:], [3]uint32{op, inkRB, inkG})
		}
	}
}

// blendRow lays the colour on row by mask: how is the opacity and the
// colour's red and blue, then green, in their places.
func blendRow(row []canvas.Color, mask []uint8, how [3]uint32) {
	op, inkRB, inkG := how[0], how[1], how[2]
	for x, m := range mask {
		if m == 0 {
			continue
		}
		a := uint32(m) * op >> 8
		inv := 255 - a
		p := uint32(row[x])
		rb := ((p&0xFF00FF)*inv + inkRB*a + 0x800080) >> 8 & 0xFF00FF
		g := ((p&0xFF00)*inv + inkG*a + 0x8000) >> 8 & 0xFF00
		row[x] = canvas.Color(p&0xFF000000 | rb | g)
	}
}
