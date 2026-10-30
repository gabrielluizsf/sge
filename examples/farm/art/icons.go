package art

// Icons are the HUD's pictures, Size square: the coin, the seed packets,
// the watering can, the axe, the harvest and the wood, and the bubble
// over one who asks for something, and the hearts, full and empty, of how
// fond someone is.
func Icons() map[string]string {
	m := map[string]string{
		"coin":        coinIcon(),
		"carrotSeeds": packetIcon(orange, leafBase),
		"tomatoSeeds": packetIcon(red, leafBase),
		"can":         canIcon(),
		"carrot":      carrotIcon(),
		"tomato":      tomatoIcon(),
		"wood":        woodIcon(),
		"bubble":      bubbleIcon(),
		"heart":       heartIcon(true),
		"heart_empty": heartIcon(false),
	}
	for name, head := range Axes {
		m["axe_"+name] = axeIcon(head)
	}
	return m
}

func coinIcon() string {
	d := newDoc(Size, Size)
	gold := d.linear(grad{0.2, 0, 0.8, 1, []stop{{0, "#fff1a6"}, {0.5, yellow}, {1, "#d69a1e"}}})
	d.add(`<circle cx="8" cy="8" r="6.2" fill="%s" stroke="%s" stroke-width="0.8"/>`, gold, line)
	d.add(`<circle cx="8" cy="8" r="4.2" fill="none" stroke="#d69a1e" stroke-width="0.6"/>`)
	d.add(`<path d="M8 5 v6 M6.6 6.4 q1.4 -1 2.8 0 q-2.8 1.6 0 3.2 q-1.4 1 -2.8 0" stroke="#b07812" stroke-width="0.7" fill="none" stroke-linecap="round"/>`)
	d.add(`<ellipse cx="5.6" cy="5.2" rx="1.4" ry="0.8" fill="%s" opacity="0.8"/>`, white)
	return d.String()
}

// packetIcon is a packet of seeds with a picture of the crop on it.
func packetIcon(fruit, leaves string) string {
	d := newDoc(Size, Size)
	d.add(`<path d="M3 3.5 L13 3.5 L12.4 14 L3.6 14 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, d.down(shirtLight, shirtShade), line)
	d.add(`<path d="M3 3.5 l1 -1.5 l1 1.5 l1 -1.5 l1 1.5 l1 -1.5 l1 1.5 l1 -1.5 l1 1.5 l1 -1.5 l1 1.5" fill="%s" stroke="%s" stroke-width="0.5"/>`, shirtLight, line)
	d.add(`<circle cx="8" cy="9.6" r="2.6" fill="%s" stroke="%s" stroke-width="0.5"/>`, fruit, line)
	d.add(`<path d="M8 7 q-1.2 -1.6 -2 -1.2 M8 7 q1.2 -1.6 2 -1.2" stroke="%s" stroke-width="0.8" fill="none" stroke-linecap="round"/>`, leaves)
	d.add(`<circle cx="7.1" cy="8.7" r="0.7" fill="%s" opacity="0.7"/>`, white)
	return d.String()
}

func canIcon() string {
	d := newDoc(Size, Size)
	body := d.down(metal, metalShade)
	d.add(`<path d="M10.5 7.6 L15 4.6" stroke="%s" stroke-width="1.6" stroke-linecap="round"/>`, line)
	d.add(`<path d="M10.5 7.6 L15 4.6" stroke="%s" stroke-width="0.8" stroke-linecap="round"/>`, metal)
	d.add(`<path d="M5 5.4 q1.6 -3 3.2 0" stroke="%s" stroke-width="1.2" fill="none"/>`, line)
	d.add(`<rect x="2.5" y="6" width="8.4" height="7.6" rx="1.6" fill="%s" stroke="%s" stroke-width="0.8"/>`, body, line)
	d.add(`<rect x="3.6" y="7" width="2" height="5" rx="1" fill="%s" opacity="0.8"/>`, white)
	d.add(`<circle cx="15.2" cy="4.4" r="0.5" fill="%s"/><circle cx="15.4" cy="6.2" r="0.4" fill="%s"/>`, waterLight, waterLight)
	return d.String()
}

func carrotIcon() string {
	d := newDoc(Size, Size)
	d.add(`<path d="M9.5 5.5 Q12.8 7.4 11.2 8.6 L4.2 14.2 Q3.4 14.6 3.8 13.6 L8.6 6.2 Q9 5.2 9.5 5.5 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, d.down(orange, orangeShade), line)
	d.add(`<path d="M7 9.4 l1.2 0.6 M5.6 11.6 l1 0.6" stroke="%s" stroke-width="0.5"/>`, orangeShade)
	d.add(`<path d="M10.4 5.6 q0 -3 1.6 -3.6 M10.8 6 q2 -2 4 -1.6 M10.2 5.4 q-1.4 -2 -0.6 -3.6" stroke="%s" stroke-width="1.1" fill="none" stroke-linecap="round"/>`, leafBase)
	return d.String()
}

func tomatoIcon() string {
	d := newDoc(Size, Size)
	d.add(`<ellipse cx="8" cy="9.4" rx="5.8" ry="5" fill="%s" stroke="%s" stroke-width="0.8"/>`, d.down(red, redShade), line)
	d.add(`<path d="M5.4 5.2 L8 6.6 L10.6 5.2 L9.4 7.4 L11.6 8 L8 8 L4.4 8 L6.6 7.4 Z" fill="%s" stroke="%s" stroke-width="0.4" stroke-linejoin="round"/>`, leafBase, line)
	d.add(`<ellipse cx="5.6" cy="9" rx="1.2" ry="1.6" fill="%s" opacity="0.6"/>`, white)
	return d.String()
}
