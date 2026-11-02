package main

import "github.com/gabrielluizsf/antui/canvas"

// shrink is src made to fit w by h pixels, its shape kept, each pixel the
// average of the ones it covers, weighted by their opacity so the edges
// of what is drawn do not darken. A picture that fits already is as it was.
func shrink(src *canvas.Canvas, w, h int) *canvas.Canvas {
	k := min(float64(w)/float64(src.Width), float64(h)/float64(src.Height), 1)
	ow, oh := max(1, int(float64(src.Width)*k+0.5)), max(1, int(float64(src.Height)*k+0.5))
	if ow == src.Width && oh == src.Height {
		return src
	}
	out, err := canvas.NewCanvas(ow, oh)
	if err != nil {
		panic(err)
	}
	for y := range oh {
		y0, y1 := y*src.Height/oh, max(y*src.Height/oh+1, (y+1)*src.Height/oh)
		for x := range ow {
			x0, x1 := x*src.Width/ow, max(x*src.Width/ow+1, (x+1)*src.Width/ow)
			var r, g, b, a, n int
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					c := src.At(sx, sy)
					al := int(c.A())
					r, g, b, a, n = r+int(c.R())*al, g+int(c.G())*al, b+int(c.B())*al, a+al, n+1
				}
			}
			if a == 0 {
				continue
			}
			out.Put(x, y, canvas.RGBA(canvas.RGBAParams{R: uint8(r / a), G: uint8(g / a), B: uint8(b / a), A: uint8(a / n)}))
		}
	}
	return out
}
