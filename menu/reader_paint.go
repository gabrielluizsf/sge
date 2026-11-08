package menu

import "github.com/gabrielluizsf/antui/canvas"

// The reader's panel, in pixels of art: its width, and how tall it is at
// most.
const (
	readerW = 300
	readerH = 190
)

// Paint draws the reader in the middle of cv.
func (r *Reader) Paint(cv *canvas.Canvas, art Art) {
	s := r.Scale
	w, h := min(readerW, cv.Width/s-8), min(readerH, cv.Height/s-8)
	if r.words == nil {
		r.words = map[string]int{}
	}
	r.wrap((w-2*edge-6)*s, Widths(art, r.words))
	x, y := (cv.Width-w*s)/2, (cv.Height-h*s)/2
	cv.Blit(x, y, art.Panel(w, h, false))
	t := art.Text(r.Title, 'y')
	cv.Blit(x+(w*s-t.Width)/2, y+edge*s+(titleRowH*s-t.Height)/2-2*s, t)
	const rowH = 11
	r.shown = max(1, (h-titleRowH-hintRowH-2*edge)/rowH)
	r.Top = r.clamp(r.Top)
	for i := range r.shown {
		if r.Top+i >= len(r.wrapped) {
			break
		}
		if line := r.wrapped[r.Top+i]; line != "" {
			cv.Blit(x+edge*s, y+(edge+titleRowH+i*rowH)*s, art.Text(line, 'k'))
		}
	}
	hint := art.Text("Setas/Page: rolar   Esc: fechar", 'k')
	cv.Blit(x+(w*s-hint.Width)/2, y+(h-hintRowH/2)*s-hint.Height/2-edge*s/2, hint)
	if n := len(r.wrapped); n > r.shown {
		bar := r.shown * rowH * s
		bh := max(4, bar*r.shown/n)
		by := y + (edge+titleRowH)*s + (bar-bh)*r.Top/(n-r.shown)
		cv.Blit(x+(w-edge/2-1)*s, by, art.Panel(2, bh/s, true))
	}
}
