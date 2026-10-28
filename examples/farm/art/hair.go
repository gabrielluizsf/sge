package art

// hairBehind is what of the hair of look l falls behind the body, drawn
// first: long hair down the back, a curly cloud, a ponytail; seen from
// the side when side is set, from the front otherwise.
func hairBehind(d *doc, l Look, side bool) {
	if l.Style != Long && l.Style != Curly && !(l.Style == Ponytail && side) {
		return
	}
	paint := d.down(l.Hair[1], l.Hair[2])
	x := 12.0
	if side {
		x = 10.6
	}
	switch l.Style {
	case Long:
		d.add(`<path d="M%.2f 8 Q%.2f 6 %.2f 8 L%.2f 19.6 Q%.2f 21 %.2f 19.6 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`,
			x-5, x, x+5, x+5.4, x, x-5.4, paint, line)
	case Curly:
		curls(d, pt{x, 9.6}, 6.4, paint)
	case Ponytail:
		if side {
			d.add(`<path d="M8.6 9.4 Q5.4 12 6.4 17.4 Q7.6 15 9.4 12 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, paint, line)
		}
	}
}

// curls is a cloud of curls, a ring of round tufts r around c, in paint.
func curls(d *doc, c pt, r float64, paint string) {
	d.add(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s" stroke="%s" stroke-width="0.6"/>`, c.x, c.y, r, paint, line)
	for i := range 10 {
		a := float64(i) / 10
		x, y := c.x+r*cosTurn(a), c.y+r*sinTurn(a)
		d.add(`<circle cx="%.2f" cy="%.2f" r="1.6" fill="%s" stroke="%s" stroke-width="0.5"/>`, x, y, paint, line)
	}
	d.add(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, c.x, c.y, r-0.4, paint)
}

// hairSides are the locks of hair beside the face, seen from the front.
func hairSides(d *doc, l Look) {
	switch l.Style {
	case Buzz, Curly:
		return
	case Long:
		d.add(`<path d="M7.6 9.6 q-1 5 -0.4 9 l2 -0.6 l0.4 -6 Z M16.4 9.6 q1 5 0.4 9 l-2 -0.6 l-0.4 -6 Z" fill="%s" stroke="%s" stroke-width="0.5"/>`, l.Hair[1], line)
		return
	}
	d.add(`<path d="M7.6 9.6 q-0.6 4 0.8 6 l1.2 -3.4 Z M16.4 9.6 q0.6 4 -0.8 6 l-1.2 -3.4 Z" fill="%s" stroke="%s" stroke-width="0.5"/>`, l.Hair[1], line)
}

// hairTop is the hair over the head seen from the front: its crown where
// no hat hides it, then the fringe of the style, and a bun on top.
func hairTop(d *doc, l Look) {
	switch l.Style {
	case Curly:
		for _, x := range []float64{8.4, 10.6, 13.4, 15.6} {
			d.add(`<circle cx="%.2f" cy="7.6" r="1.7" fill="%s" stroke="%s" stroke-width="0.5"/>`, x, d.down(l.Hair[0], l.Hair[1]), line)
		}
		return
	case Buzz:
		d.add(`<path d="M7.8 10 Q7.6 5.6 12 5.6 Q16.4 5.6 16.2 10 Q12 7.4 7.8 10 Z" fill="%s" stroke="%s" stroke-width="0.5" opacity="0.95"/>`, l.Hair[1], line)
		return
	}
	if !l.Hat {
		d.add(`<path d="M7.4 10 Q7.2 5 12 4.8 Q16.8 5 16.6 10 Z" fill="%s" stroke="%s" stroke-width="0.6"/>`, d.down(l.Hair[0], l.Hair[1]), line)
		if l.Style == Bun {
			d.add(`<circle cx="12" cy="4" r="2.4" fill="%s" stroke="%s" stroke-width="0.6"/>`, d.down(l.Hair[0], l.Hair[2]), line)
		}
	}
	bangs(d, l, 7.6)
}

// hairBack is the back of the head, all hair, as the style lies.
func hairBack(d *doc, l Look) {
	paint := d.down(l.Hair[0], l.Hair[2])
	switch l.Style {
	case Curly:
		curls(d, pt{12, 9.8}, 6.2, paint)
		return
	case Long:
		d.add(`<path d="M7.2 10 Q7 5.4 12 5.4 Q17 5.4 16.8 10 L17.2 19.4 Q12 20.6 6.8 19.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, paint, line)
		return
	}
	d.add(`<ellipse cx="12" cy="10.6" rx="4.6" ry="5" fill="%s" stroke="%s" stroke-width="0.7"/>`, paint, line)
	if l.Style != Buzz {
		d.add(`<path d="M9.4 13.4 q1 1.8 2.6 2 q1.6 -0.2 2.6 -2 M10.4 9 q0.6 3 0.2 5 M13.6 9 q-0.6 3 -0.2 5" stroke="%s" stroke-width="0.5" fill="none"/>`, l.Hair[2])
	}
	switch {
	case l.Style == Ponytail:
		d.add(`<path d="M10.8 12 Q12 11 13.2 12 L13 18.6 Q12 19.6 11 18.6 Z" fill="%s" stroke="%s" stroke-width="0.6"/>`, paint, line)
	case l.Style == Bun && !l.Hat:
		d.add(`<circle cx="12" cy="6.4" r="2.4" fill="%s" stroke="%s" stroke-width="0.6"/>`, paint, line)
	}
}

// hairSide is the hair seen from the side, over the head in profile.
func hairSide(d *doc, l Look) {
	paint := d.down(l.Hair[0], l.Hair[1])
	switch l.Style {
	case Curly:
		for _, c := range []pt{{9.6, 7.4}, {12.4, 6.4}, {14.8, 7.6}} {
			d.add(`<circle cx="%.2f" cy="%.2f" r="1.8" fill="%s" stroke="%s" stroke-width="0.5"/>`, c.x, c.y, paint, line)
		}
		return
	case Buzz:
		d.add(`<path d="M15.4 8.8 Q13 5.8 10 7.2 Q8.6 9.4 9.6 12.6 L11.2 10.6 Q13.4 9.4 15.4 8.8 Z" fill="%s" stroke="%s" stroke-width="0.5"/>`, l.Hair[1], line)
		return
	}
	if !l.Hat {
		d.add(`<path d="M8 10.4 Q7.6 5 12.4 5 Q16.2 5.2 16 9.2 Z" fill="%s" stroke="%s" stroke-width="0.6"/>`, paint, line)
		if l.Style == Bun {
			d.add(`<circle cx="9" cy="5.6" r="2.3" fill="%s" stroke="%s" stroke-width="0.6"/>`, d.down(l.Hair[0], l.Hair[2]), line)
		}
	}
	d.add(`<path d="M13.6 7.6 Q8.4 6.4 7.8 10.6 Q7.6 13.6 9.4 15.6 L10.8 13 L11.6 14.4 L12.2 11.4 Q13.6 10.6 14.8 9.6 L15.8 9.8 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, paint, line)
}
