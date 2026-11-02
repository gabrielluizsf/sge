package backpack

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// painter is a pack being drawn: onto cv, with art.
type painter struct {
	p   *Pack
	cv  *canvas.Canvas
	art Art
}

// Paint draws the hotbar and, if the screen is open, the grid over it, the
// item held following the mouse or the cursor, and a tooltip. It places
// the pack on cv first.
func (p *Pack) Paint(cv *canvas.Canvas, art Art) {
	hint := art.Text(p.Hint, 'k')
	p.Place(View{W: cv.Width, H: cv.Height, Scale: p.Scale, MinW: hint.Width + 2*padding*p.Scale})
	d := painter{p: p, cv: cv, art: art}
	if p.Open {
		d.screen()
	}
	for i := range p.Hotbar {
		d.slot(Spot{Bar, i})
	}
	if p.Open {
		d.heldItem()
		d.tooltip()
	}
}

// screen draws the panel, its title and hint, and the grid's slots.
func (d painter) screen() {
	p, s := d.p, d.p.Scale
	r := p.rects.panel
	d.cv.Blit(r.X, r.Y, d.art.Panel(r.W/s, r.H/s, false))
	title := d.art.Text(p.Title, 'k')
	d.cv.Blit(r.X+padding*s, r.Y+(titleH*s-title.Height)/2+s, title)
	hint := d.art.Text(p.Hint, 'k')
	d.cv.Blit(r.X+(r.W-hint.Width)/2, r.Y+r.H-(hintH*s+hint.Height)/2-s, hint)
	for i := range p.rects.grid {
		d.slot(Spot{Grid, i})
	}
}

// slot draws slot sp: its panel, gold under the cursor, or for the hotbar
// slot in hand when the screen is closed; its item's picture and count;
// and a hotbar slot's number.
func (d painter) slot(sp Spot) {
	p := d.p
	r := p.rectOf(sp)
	lit := p.Open && p.cursor == sp || !p.Open && sp.Area == Bar && sp.Index == p.Active
	d.cv.Blit(r.X, r.Y, d.art.Panel(SlotSize, SlotSize, lit))
	if sp.Area == Bar {
		k := d.art.Text(fmt.Sprint(sp.Index+1), 'w')
		d.cv.Blit(r.X+(r.W-k.Width)/2, r.Y-k.Height, k)
	}
	id := p.itemAt(sp)
	if id == "" || p.held == id && p.dragging {
		return
	}
	if it, ok := p.Source.Find(id); ok {
		d.item(it, r)
	}
}

// item draws it in the slot at r: its picture and, below 0 aside, its
// count at the corner.
func (d painter) item(it Item, r Rect) {
	s := d.p.Scale
	d.cv.Blit(r.X+2*s, r.Y+2*s, d.art.Icon(it.Icon))
	if it.Count >= 0 {
		t := d.art.Text(fmt.Sprint(it.Count), 'w')
		d.cv.Blit(r.X+r.W-t.Width+2, r.Y+r.H-t.Height+4, t)
	}
}
