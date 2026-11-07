package art

// beastLook is how a four-legged animal is drawn from the side, facing
// right, in art units: its body, how long its legs are, its colours.
type beastLook struct {
	w, h                      float64 // the picture
	bodyW, bodyH, bodyY       float64 // the body's ellipse, its middle at bodyY above the feet
	legH, legW                float64
	head, ear                 float64 // the head's size, and the ear's
	snout, coat, patch, shade string
	tail, horn                string // colours, "" for none
	spots                     bool
}

var beasts = map[string]beastLook{
	"cow":   {w: 34, h: 26, bodyW: 11, bodyH: 6.5, bodyY: 11.5, legH: 8, legW: 2.2, head: 4.2, ear: 1.6, snout: pink, coat: white, patch: "#3b3436", shade: "#d9d4d6", tail: "#3b3436", horn: "#f4e9c8", spots: true},
	"sheep": {w: 26, h: 20, bodyW: 8, bodyH: 5.5, bodyY: 9, legH: 6, legW: 1.5, head: 3.2, ear: 1.4, snout: "#3a3036", coat: "#f4f0e6", patch: "#f4f0e6", shade: "#cfc7b6", tail: "", horn: ""},
	"horse": {w: 40, h: 34, bodyW: 13, bodyH: 6.5, bodyY: 16, legH: 12, legW: 2.4, head: 5, ear: 1.6, snout: "#3a2a22", coat: "#a5612f", patch: "#a5612f", shade: "#7b4421", tail: "#e9d49a", horn: ""},
	"wolf":  {w: 30, h: 20, bodyW: 9.5, bodyH: 4.6, bodyY: 9.5, legH: 7, legW: 1.7, head: 3.8, ear: 1.7, snout: "#2a2428", coat: "#6f7076", patch: "#6f7076", shade: "#4b4c52", tail: "#55565c", horn: ""},
}

// coats are the colours a kind may come in, apart from its first: the
// horse's second is grey.
var coats = map[string]string{"horse": "#b8bcc4"}

// Animal is the SVG of an animal of kind, "cow", "sheep", "horse" or
// "wolf", facing right, and how wide and tall its picture is. The feet are
// at the bottom, in the middle. frame is a step of its walk, 0 standing;
// down lowers its head, as when it grazes, and lying turns it onto its
// side, a dead wolf's.
func Animal(kind string, coat, frame int, down, lying bool) (svg string, w, h int) {
	b, ok := beasts[kind]
	if !ok {
		b = beasts["cow"]
	}
	if c, ok := coats[kind]; ok && coat > 0 {
		b.coat, b.patch, b.shade = c, c, "#8d929c"
	}
	d := newDoc(b.w, b.h)
	cx := b.w/2 - 1
	by := b.h - b.bodyY - 1.5
	if lying {
		by = b.h - b.bodyH - 1
	}
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s"/>`, cx, b.h-1.2, b.bodyW*0.95, 1.8, "#1f3a22")
	legs := [][2]float64{{-0.62, 0}, {-0.45, 1}, {0.42, 1}, {0.6, 0}}
	for i, l := range legs {
		swing := 0.0
		if frame > 0 {
			swing = []float64{-2, 0, 2, 0}[(frame+int(l[1])*2+i/2)%4] * b.legW / 2.2
		}
		x := cx + l[0]*b.bodyW
		shade := b.shade
		if i == 1 || i == 2 {
			shade = b.coat
		}
		if lying {
			d.add(`<rect x="%g" y="%g" width="%g" height="%g" rx="0.6" fill="%s"/>`, x-0.5, by+b.bodyH*0.3+float64(i)*0.8, b.legH*0.8, b.legW, shade)
			continue
		}
		d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`, x, by+b.bodyH*0.5, x+swing, b.h-1.5, shade, b.legW)
	}
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" stroke="%s" stroke-width="0.5"/>`, cx, by, b.bodyW, b.bodyH, b.coat, line)
	if b.spots {
		d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s"/>`, cx-2, by-1, 3.6, 2.4, b.patch)
		d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s"/>`, cx+3.5, by+1.5, 2.4, 1.8, b.patch)
	}
	if b.tail != "" {
		d.add(`<path d="M%g %g Q%g %g %g %g" stroke="%s" stroke-width="1.6" fill="none" stroke-linecap="round"/>`, cx-b.bodyW, by-b.bodyH*0.3, cx-b.bodyW-2.5, by+b.bodyH*0.2, cx-b.bodyW-1.5, by+b.bodyH*1.1, b.tail)
	}
	hx, hy := cx+b.bodyW+b.head*0.3, by-b.bodyH*0.8
	if down {
		hx, hy = cx+b.bodyW+b.head*0.5, by+b.bodyH*0.5
	}
	d.add(`<path d="M%g %g L%g %g" stroke="%s" stroke-width="%g" stroke-linecap="round"/>`, cx+b.bodyW*0.8, by-b.bodyH*0.2, hx, hy, b.coat, b.head*0.8)
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s" stroke="%s" stroke-width="0.5"/>`, hx+b.head*0.4, hy+b.head*0.2, b.head*0.9, b.head*0.65, b.coat, line)
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s"/>`, hx+b.head*1.1, hy+b.head*0.4, b.head*0.45, b.head*0.4, b.snout)
	d.add(`<ellipse cx="%g" cy="%g" rx="%g" ry="%g" fill="%s"/>`, hx-b.head*0.1, hy-b.head*0.45, b.ear*0.6, b.ear, b.shade)
	if b.horn != "" {
		d.add(`<path d="M%g %g l-0.5 -2" stroke="%s" stroke-width="1" stroke-linecap="round"/>`, hx+b.head*0.2, hy-b.head*0.5, b.horn)
	}
	d.add(`<circle cx="%g" cy="%g" r="0.7" fill="%s"/>`, hx+b.head*0.6, hy, eye)
	return d.String(), int(b.w), int(b.h)
}
