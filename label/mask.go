package label

import "github.com/gabrielluizsf/antui/canvas"

// mask is how much of each pixel the glyphs of s cover, from 0 to 255, w
// pixels a row: drawn white on black, so each pixel's grey is its
// coverage, whatever the canvas does with alpha.
func mask(face *canvas.Face, s string) (cover []uint8, w, h int) {
	w, h = max(1, face.Width(s)), face.Height()
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		panic(err)
	}
	cv.Clear(canvas.RGB(0, 0, 0))
	face.Draw(canvas.FaceDrawParams{Canvas: cv, Pos: canvas.Point{Y: float64(face.Ascent())}, Text: s, Color: canvas.RGB(255, 255, 255)})
	cover = make([]uint8, w*h)
	for y := range h {
		for x := range w {
			cover[y*w+x] = cv.At(x, y).G()
		}
	}
	return cover, w, h
}
