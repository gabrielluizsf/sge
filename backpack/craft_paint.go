package backpack

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// tabs draws the tabs' names on the panel's top, the one shown lit.
func (d painter) tabs() {
	p, s := d.p, d.p.Scale
	p.rects.tabs = p.rects.tabs[:0]
	x := p.rects.panel.X + padding*s
	for i, name := range p.Tabs {
		t := d.art.Text(name, 'k')
		w := t.Width/s + 12
		r := Rect{x, p.rects.panel.Y + s, w * s, tabH * s}
		p.rects.tabs = append(p.rects.tabs, r)
		d.cv.Blit(r.X, r.Y, d.art.Panel(w, tabH, i == p.Tab))
		d.cv.Blit(r.X+(r.W-t.Width)/2, r.Y+(r.H-t.Height)/2+s, t)
		x += r.W + 2*s
	}
}

// recipeRows draws the rows of the Craft tab that show: each recipe's
// picture and name, and its materials at the right, what is carried over
// what is wanted, greyed when one is lacking.
func (d painter) recipeRows() {
	p := d.p
	rs := p.recipes()
	for k, r := range p.rects.rows {
		i := p.craft.top + k
		lit := i == p.craft.cursor
		d.cv.Blit(r.X, r.Y, d.art.Panel(r.W/p.Scale, r.H/p.Scale, lit))
		if i >= len(rs) {
			continue
		}
		d.recipe(rs[i], r)
	}
}

// recipe draws recipe rc on the row at r.
func (d painter) recipe(rc Recipe, r Rect) {
	s := d.p.Scale
	can := rc.Can()
	icon := d.art.Icon(rc.Icon)
	if !can {
		icon = Greyed(icon)
	}
	d.cv.Blit(r.X+2*s, r.Y+2*s, icon)
	name := d.art.Text(rc.Name, 'k')
	if !can {
		name = Greyed(name)
	}
	d.cv.Blit(r.X+22*s, r.Y+(r.H-name.Height)/2+s, name)
	x := r.X + r.W - 2*s
	for i := len(rc.Needs) - 1; i >= 0; i-- {
		n := rc.Needs[i]
		ink := rune('g')
		if n.Have < n.Want {
			ink = 'r'
		}
		t := d.art.Text(fmt.Sprintf("%d/%d", n.Have, n.Want), ink)
		x -= t.Width
		d.cv.Blit(x, r.Y+(r.H-t.Height)/2+s, t)
		ic := d.art.Icon(n.Icon)
		x -= ic.Width + s
		d.cv.Blit(x, r.Y+(r.H-ic.Height)/2, ic)
		x -= 3 * s
	}
}

// Greyed is a copy of cv in greys, as dim as it was bright and as seen
// through as it was.
func Greyed(cv *canvas.Canvas) *canvas.Canvas {
	out, err := canvas.NewCanvas(cv.Width, cv.Height)
	if err != nil {
		return cv
	}
	for y := range cv.Height {
		for x := range cv.Width {
			c := cv.At(x, y)
			l := (uint32(c.R())*3 + uint32(c.G())*6 + uint32(c.B())) / 10
			l = 70 + l*2/5
			out.Pixels[y*out.Stride+x] = canvas.RGBA(canvas.RGBAParams{R: uint8(l), G: uint8(l), B: uint8(l), A: c.A()})
		}
	}
	return out
}
