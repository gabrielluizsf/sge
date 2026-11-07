package art

// heart is the outline of a heart, Size square.
const heart = `M8 14 C3 10.4 1.4 8 1.4 5.6 Q1.4 2.4 4.6 2.2 Q6.8 2.2 8 4.4 Q9.2 2.2 11.4 2.2 Q14.6 2.4 14.6 5.6 C14.6 8 13 10.4 8 14 Z`

// heartIcon is a heart, as games show how fond someone is: pink to red
// and shining when full, a pale outline of one when empty.
func heartIcon(full bool) string {
	d := newDoc(Size, Size)
	if !full {
		d.add(`<path d="%s" fill="%s" fill-opacity="0.45" stroke="%s" stroke-width="0.9" stroke-linejoin="round"/>`, heart, white, line)
		return d.String()
	}
	d.add(`<path d="%s" fill="%s" stroke="%s" stroke-width="0.9" stroke-linejoin="round"/>`, heart, d.down(pink, red), line)
	d.add(`<ellipse cx="4.8" cy="5" rx="1.4" ry="0.9" fill="%s" opacity="0.8"/>`, white)
	return d.String()
}
