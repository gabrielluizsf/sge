package art

// woodIcon is a bundle of logs, their cut ends to the viewer.
func woodIcon() string {
	d := newDoc(Size, Size)
	bark := d.down(barkLight, barkShade)
	for _, c := range [][2]float64{{5, 10.6}, {11, 10.6}, {8, 6}} {
		d.add(`<rect x="%.2f" y="%.2f" width="6" height="5" rx="2.5" fill="%s" stroke="%s" stroke-width="0.6"/>`, c[0]-3, c[1]-2.5, bark, line)
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="1.8" ry="2.1" fill="%s" stroke="%s" stroke-width="0.4"/>`, c[0]+1.2, c[1], woodLight, line)
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="0.8" ry="1" fill="none" stroke="%s" stroke-width="0.3"/>`, c[0]+1.2, c[1], woodShade)
	}
	return d.String()
}

// Axes are the heads of the axes, light and shade, from the worst to the
// best: stone, iron, gold and diamond.
var Axes = map[string][2]string{
	"stone":   {"#b3aea6", "#77716a"},
	"iron":    {metal, metalShade},
	"gold":    {"#fff1a6", "#d69a1e"},
	"diamond": {"#c9fbff", "#3fc6d6"},
}

// axeIcon is an axe: a wooden handle and a head of what head paints.
func axeIcon(head [2]string) string {
	d := newDoc(Size, Size)
	d.add(`<path d="M4 14 L11.4 3.6" stroke="%s" stroke-width="2.2" stroke-linecap="round"/>`, line)
	d.add(`<path d="M4 14 L11.4 3.6" stroke="%s" stroke-width="1.2" stroke-linecap="round"/>`, woodBase)
	d.add(`<path d="M9.2 3.2 Q12.6 0.6 14.8 3.4 Q13.6 6.6 11.6 7.6 L9.2 5.6 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.down(head[0], head[1]), line)
	d.add(`<path d="M12.6 2 Q14.2 2.6 14.6 3.6" stroke="%s" stroke-width="0.5" fill="none"/>`, white)
	return d.String()
}

// bubbleIcon is a speech bubble with an exclamation mark: one who has
// something to ask.
func bubbleIcon() string {
	d := newDoc(Size, Size)
	d.add(`<path d="M2 3.6 Q2 1.6 4 1.6 L12 1.6 Q14 1.6 14 3.6 L14 9.4 Q14 11.4 12 11.4 L8.6 11.4 L6 14.4 L6.2 11.4 L4 11.4 Q2 11.4 2 9.4 Z" fill="%s" stroke="%s" stroke-width="0.8" stroke-linejoin="round"/>`, white, line)
	d.add(`<path d="M8 3.6 L8 7.6" stroke="%s" stroke-width="1.8" stroke-linecap="round"/>`, orangeShade)
	d.add(`<circle cx="8" cy="9.6" r="0.95" fill="%s"/>`, orangeShade)
	return d.String()
}

// anvilTile is the smith's anvil on a block of wood, standing on its
// cell.
func anvilTile() string {
	d := newDoc(IsoW, IsoH)
	shadow(d, 7, 3)
	d.origin = pt{16, 32}
	d.box([4]float64{5, 5, 11, 11}, [2]float64{0, 5}, [3]string{woodLight, woodBase, woodShade})
	d.box([4]float64{6.5, 6.5, 9.5, 9.5}, [2]float64{5, 8}, [3]string{slateLight, slateBase, slateShade})
	d.box([4]float64{3, 5.5, 13, 10.5}, [2]float64{8, 11}, [3]string{slateLight, slateBase, slateShade})
	d.add(`<path d="M%s L%s L%s Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, d.at(13, 6.5, 11), d.at(13, 9.5, 11), d.at(16, 8, 10.4), slateLight, line)
	d.add(`<path d="M%s L%s" stroke="%s" stroke-width="0.6" stroke-linecap="round"/>`, d.at(4, 6.5, 11), d.at(12, 6.5, 11), white)
	return d.String()
}
