package art

// pickaxeIcon is a pickaxe: a wooden handle and, across its end, an iron
// head of two points curving down.
func pickaxeIcon() string {
	d := newDoc(Size, Size)
	d.add(`<path d="M3.6 14 L11 5.4" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>`, line)
	d.add(`<path d="M3.6 14 L11 5.4" stroke="%s" stroke-width="1.2" stroke-linecap="round"/>`, woodBase)
	d.add(`<path d="M5.8 4.6 Q9.4 0.8 14.6 4.4 Q10.8 3.8 9.6 5.6 Q8 4.2 5.8 4.6 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.down(metal, metalShade), line)
	d.add(`<path d="M8.6 3.2 Q11 2.4 12.8 3.4" stroke="%s" stroke-width="0.5" fill="none"/>`, white)
	return d.String()
}
