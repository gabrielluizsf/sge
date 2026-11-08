package label

import "github.com/gabrielluizsf/antui/canvas"

// New is s in face as st draws it: a picture the size of the text and its
// outline, the ink over the outline, each pixel as opaque as the glyphs
// cover it.
func New(face *canvas.Face, s string, st Style) *canvas.Canvas {
	cover, mw, mh := mask(face, s)
	pad := max(0, st.Width)
	w, h := mw+2*pad, mh+2*pad
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		panic(err)
	}
	at := func(x, y int) float64 {
		if x < 0 || y < 0 || x >= mw || y >= mh {
			return 0
		}
		return float64(cover[y*mw+x]) / 255
	}
	for y := range h {
		for x := range w {
			ink := at(x-pad, y-pad)
			edge := 0.0
			for dy := -pad; dy <= pad && pad > 0; dy++ {
				for dx := -pad; dx <= pad; dx++ {
					if dx*dx+dy*dy <= pad*pad+1 {
						edge = max(edge, at(x-pad+dx, y-pad+dy))
					}
				}
			}
			cv.Pixels[y*cv.Stride+x] = over(colour{st.Ink, ink}, colour{st.Outline, edge})
		}
	}
	return cv
}

// colour is a colour as opaque as a coverage.
type colour struct {
	c     canvas.Color
	cover float64
}

// over is top laid over bottom, each as opaque as its colour's alpha
// times its cover: their colours mixed by how much each shows.
func over(top, bottom colour) canvas.Color {
	a := top.cover * float64(top.c.A()) / 255
	b := bottom.cover * float64(bottom.c.A()) / 255 * (1 - a)
	alpha := a + b
	if alpha <= 0 {
		return canvas.Transparent
	}
	mix := func(p, q uint8) uint8 { return uint8((float64(p)*a+float64(q)*b)/alpha + 0.5) }
	return canvas.RGBA(canvas.RGBAParams{
		R: mix(top.c.R(), bottom.c.R()), G: mix(top.c.G(), bottom.c.G()), B: mix(top.c.B(), bottom.c.B()),
		A: uint8(alpha*255 + 0.5),
	})
}
