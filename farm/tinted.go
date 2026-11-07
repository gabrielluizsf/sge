package farm

import "github.com/gabrielluizsf/antui/canvas"

// tinted is the picture draw makes with every pixel pushed toward tint and
// made translucent, k of the way opaque; it is drawn once, kept as key.
func (f *farm) tinted(key string, draw func() *canvas.Canvas, tint canvas.Color, k float64) *canvas.Canvas {
	if out, ok := f.tints[key]; ok {
		return out
	}
	cv := draw()
	out, err := canvas.NewCanvas(cv.Width, cv.Height)
	if err != nil {
		return cv
	}
	for y := range cv.Height {
		for x := range cv.Width {
			c := cv.At(x, y)
			if c.A() == 0 {
				continue
			}
			mix := func(a, b uint8) uint8 { return uint8((int(a) + int(b)) / 2) }
			out.Pixels[y*out.Stride+x] = canvas.RGBA(canvas.RGBAParams{R: mix(c.R(), tint.R()), G: mix(c.G(), tint.G()), B: mix(c.B(), tint.B()), A: uint8(float64(c.A()) * k)})
		}
	}
	if f.tints == nil {
		f.tints = map[string]*canvas.Canvas{}
	}
	f.tints[key] = out
	return out
}
