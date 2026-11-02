package backpack

import "github.com/gabrielluizsf/antui/canvas"

// heldItem draws the item taken up: under the mouse while it drags it,
// and lifted over the cursor's slot when a key took it up.
func (d painter) heldItem() {
	p := d.p
	it, ok := p.Source.Find(p.held)
	if !ok {
		return
	}
	s := p.Scale
	r := p.rectOf(p.cursor)
	r.X, r.Y = r.X-3*s, r.Y-3*s
	if p.dragging {
		r.X, r.Y = p.at.X-SlotSize*s/2, p.at.Y-SlotSize*s/2
	}
	d.cv.Blit(r.X, r.Y, d.art.Panel(SlotSize, SlotSize, true))
	d.item(it, r)
}

// tooltip draws the name, and what it is for, of the item under the
// cursor, under the screen's panel, where it covers no slot, or over it if there is no
// room, unless one is held.
func (d painter) tooltip() {
	p := d.p
	if p.held != "" {
		return
	}
	it, ok := p.Source.Find(p.itemAt(p.cursor))
	if !ok {
		return
	}
	s := p.Scale
	texts := []*canvas.Canvas{d.art.Text(it.Name, 'y')}
	for _, l := range it.Info {
		texts = append(texts, d.art.Text(l, 'k'))
	}
	w, rows := 0, 0
	for _, t := range texts {
		w = max(w, t.Width)
		rows += t.Height + 2*s
	}
	pw, ph := w/s+12, rows/s+8
	panel := p.rects.panel
	x, y := panel.X, panel.Y+panel.H+4*s
	if bar := p.rects.bar; len(bar) > 0 && y+ph*s > bar[0].Y-8*s {
		y = panel.Y - ph*s - 4*s
	}
	x, y = max(0, x), max(0, y)
	d.cv.Blit(x, y, d.art.Panel(pw, ph, false))
	at := y + 4*s
	for _, t := range texts {
		d.cv.Blit(x+6*s, at, t)
		at += t.Height + 2*s
	}
}
