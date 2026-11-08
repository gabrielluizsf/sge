package render3d

import "image"

// rgba8 is img's size and its pixels as RGBA8 bytes, top row first, not
// premultiplied.
func rgba8(img image.Image) (w, h int, px []byte) {
	b := img.Bounds()
	w, h = b.Dx(), b.Dy()
	if n, ok := img.(*image.NRGBA); ok && n.Stride == 4*w {
		return w, h, n.Pix[:4*w*h]
	}
	px = make([]byte, 0, 4*w*h)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()
			if a == 0 {
				px = append(px, 0, 0, 0, 0)
				continue
			}
			// RGBA is premultiplied: undo it.
			px = append(px, byte(r*0xff/a), byte(g*0xff/a), byte(bl*0xff/a), byte(a>>8))
		}
	}
	return w, h, px
}
