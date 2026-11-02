package dialogue

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/menu"
)

// The box's measures, in pixels of art.
const (
	boxW      = 400
	boxH      = 88
	margin    = 10
	portraitW = 56
	pad       = 6
	buttonW   = 110
	buttonH   = 12
	buttonGap = 2
	textRows  = 4
	rowH      = 12
	nameH     = 14
)

// Paint draws the box at the bottom of cv, centred, with portrait, which
// may be nil, on its left, and lays out its pages and buttons.
func (b *Box) Paint(cv *canvas.Canvas, art Art, portrait *canvas.Canvas) {
	s := b.Scale
	w := min(boxW, cv.Width/s-2*margin)
	b.frame = Rect{X: (cv.Width - w*s) / 2, Y: cv.Height - (boxH+margin)*s, W: w * s, H: boxH * s}
	f := b.frame
	textW := (w - portraitW - buttonW - 4*pad) * s
	if !b.laid {
		if b.words == nil {
			b.words = map[string]int{}
		}
		lines := menu.Wrap(b.text, textW, menu.Widths(art, b.words))
		b.pages = nil
		for i := 0; i < len(lines) || i == 0; i += textRows {
			b.pages = append(b.pages, lines[i:min(i+textRows, len(lines))])
		}
		b.laid = true
	}
	cv.Blit(f.X, f.Y, art.Panel(w, boxH, false))
	b.portrait(cv, art, portrait)
	tx := f.X + (portraitW+2*pad)*s
	name := art.Text(b.Speaker, 'y')
	cv.Blit(tx, f.Y+pad*s+(nameH*s-name.Height)/2, name)
	b.sayWords(cv, art, Point{X: tx, Y: f.Y + (pad+nameH)*s})
	b.buttonsOf(cv, art)
}

// portrait draws the speaker's portrait on a panel at the left.
func (b *Box) portrait(cv *canvas.Canvas, art Art, p *canvas.Canvas) {
	s, f := b.Scale, b.frame
	h := boxH - 2*pad
	x, y := f.X+pad*s, f.Y+pad*s
	cv.Blit(x, y, art.Panel(portraitW, h, false))
	if p == nil {
		return
	}
	saved := cv.Clip
	cv.Clip = canvas.Area{X: x, Y: y, Width: portraitW * s, Height: h * s}
	cv.Blit(x+(portraitW*s-p.Width)/2, y+h*s-p.Height, p)
	cv.Clip = saved
}

// sayWords draws the page being said, as many letters of it as have
// appeared, and a mark under it while there is more.
func (b *Box) sayWords(cv *canvas.Canvas, art Art, at Point) {
	if b.page >= len(b.pages) {
		return
	}
	left := int(b.shown)
	for i, line := range b.pages[b.page] {
		n := len([]rune(line))
		if left <= 0 || n == 0 {
			left -= n
			continue
		}
		t := art.Text(line, 'k')
		y := at.Y + i*rowH*b.Scale
		if left >= n {
			cv.Blit(at.X, y, t)
		} else {
			saved := cv.Clip
			cv.Clip = canvas.Area{X: at.X, Y: y, Width: t.Width * left / n, Height: t.Height}
			cv.Blit(at.X, y, t)
			cv.Clip = saved
		}
		left -= n
	}
	if b.Done() && !b.Last() {
		m := art.Text("»", 'y')
		f := b.frame
		cv.Blit(f.X+f.W-(buttonW+pad+4)*b.Scale-m.Width, f.Y+f.H-(pad+3)*b.Scale-m.Height, m)
	}
}

// buttonsOf draws the buttons on the right, the chosen on a gold panel,
// and the ones in sight when there are more than fit.
func (b *Box) buttonsOf(cv *canvas.Canvas, art Art) {
	s, f := b.Scale, b.frame
	fit := (boxH - 2*pad + buttonGap) / (buttonH + buttonGap)
	n := min(fit, len(b.Choices))
	b.buttons = make([]Rect, n)
	b.keep()
	x := f.X + f.W - (buttonW+pad)*s
	for i := range n {
		r := Rect{X: x, Y: f.Y + (pad+i*(buttonH+buttonGap))*s, W: buttonW * s, H: buttonH * s}
		b.buttons[i] = r
		cv.Blit(r.X, r.Y, art.Panel(buttonW, buttonH, b.first+i == b.Choice))
		t := art.Text(b.Choices[b.first+i], 'k')
		cv.Blit(r.X+(r.W-t.Width)/2, r.Y+(r.H-t.Height)/2, t)
	}
}
