package menu

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// The list's measures, in pixels of art.
const (
	lineH     = 16
	lineSubH  = 26
	titleRowH = 18
	hintRowH  = 14
	edge      = 6
)

// rowHeight is the height of a line of l: taller if any has a smaller
// line under its text.
func (l *List) rowHeight() int {
	for _, ln := range l.Lines {
		if ln.Sub != "" {
			return lineSubH
		}
	}
	return lineH
}

// Paint draws the list in the middle of cv, and places its lines for the
// mouse.
func (l *List) Paint(cv *canvas.Canvas, art Art) {
	s, rh := l.Scale, l.rowHeight()
	w := l.MinW
	for _, ln := range l.Lines {
		wide := max(art.Text(ln.Text, 'k').Width, art.Text(ln.Sub, 'k').Width)
		w = max(w, wide/s+buttonW(art, ln, s)+2*edge+8)
	}
	w = min(w, cv.Width*9/10/s)
	l.shown = max(0, min(len(l.Lines), (cv.Height*8/10/s-titleRowH-hintRowH-2*edge)/rh))
	l.keep()
	h := titleRowH + l.shown*rh + hintRowH + 2*edge
	x, y := (cv.Width-w*s)/2, (cv.Height-h*s)/2
	cv.Blit(x, y, art.Panel(w, h, false))
	t := art.Text(l.Title, 'y')
	cv.Blit(x+(w*s-t.Width)/2, y+edge*s+(titleRowH*s-t.Height)/2-2*s, t)
	l.rows = l.rows[:0]
	for i := range l.shown {
		r := Rect{X: x + edge*s, Y: y + (edge+titleRowH+i*rh)*s, W: (w - 2*edge) * s, H: rh * s}
		l.rows = append(l.rows, r)
		l.line(cv, art, row{l.Lines[l.Top+i], r, l.Top+i == l.Choice})
	}
	hint := art.Text(l.Hint, 'k')
	cv.Blit(x+(w*s-hint.Width)/2, y+(h-hintRowH/2)*s-hint.Height/2-edge*s/2, hint)
	l.scrollbar(cv, art, Rect{X: x + (w-edge/2-1)*s, Y: y + (edge+titleRowH)*s, W: 2 * s, H: l.shown * rh * s})
}

// buttonW is the width, in pixels of art, of the button of line ln, or 0.
func buttonW(art Art, ln Line, s int) int {
	if ln.Button == "" {
		return 0
	}
	return art.Text(ln.Button, 'k').Width/s + 10
}

// row is a line of a list in the row at a part of the screen, chosen or not.
type row struct {
	ln     Line
	at     Rect
	chosen bool
}

// line draws ln in the row at r, on a gold bar when chosen.
func (l *List) line(cv *canvas.Canvas, art Art, rw row) {
	s, ln, r, chosen := l.Scale, rw.ln, rw.at, rw.chosen
	if chosen {
		cv.Blit(r.X, r.Y, art.Panel(r.W/s, r.H/s, true))
	}
	room := r.W - 12*s - buttonW(art, ln, s)*s
	t := fit(art, ln.Text, room)
	ty := r.Y + (r.H-t.Height)/2
	if ln.Sub != "" {
		ty = r.Y + 3*s
		sub := fit(art, ln.Sub, room)
		cv.Blit(r.X+6*s, r.Y+r.H-sub.Height-2*s, sub)
	}
	cv.Blit(r.X+6*s, ty, t)
	if ln.Button != "" {
		b := art.Text(ln.Button, 'y')
		bw, bh := b.Width/s+10, b.Height/s+4
		bx, by := r.X+r.W-(bw+4)*s, r.Y+(r.H-bh*s)/2
		cv.Blit(bx, by, art.Panel(bw, bh, chosen))
		cv.Blit(bx+5*s, by+(bh*s-b.Height)/2, b)
	}
}

// fit is the text s, cut short with an ellipsis if it is wider than w
// pixels.
func fit(art Art, s string, w int) *canvas.Canvas {
	t := art.Text(s, 'k')
	for r := []rune(s); t.Width > w && len(r) > 1; {
		r = r[:len(r)-1]
		t = art.Text(strings.TrimRight(string(r), " ")+"…", 'k')
	}
	return t
}

// scrollbar draws a thin bar at r, as long as the lines in sight are of
// all, where they are, if some are out of sight.
func (l *List) scrollbar(cv *canvas.Canvas, art Art, r Rect) {
	n := len(l.Lines)
	if n <= l.shown {
		return
	}
	h := max(4, r.H*l.shown/n)
	cv.Blit(r.X, r.Y+(r.H-h)*l.Top/(n-l.shown), art.Panel(r.W/l.Scale, h/l.Scale, true))
}
