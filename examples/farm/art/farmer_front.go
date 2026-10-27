package art

// pose is how the farmer stands in a frame: how far each leg is lifted
// mid-stride, how far the arms swing, and whether they hold the can out.
type pose struct {
	lift  [2]float64 // left leg, right leg, as the viewer sees them
	swing float64
	can   bool
}

// front is the farmer facing the viewer, or with back true facing away.
func front(p pose, back bool) string {
	d := newDoc(FarmerW, FarmerH)
	d.add(`<ellipse cx="12" cy="39.2" rx="6" ry="1.3" fill="%s" opacity="0.3"/>`, shadowInk)
	// Legs and boots, each lifted as the stride says.
	for i, x := range []float64{9.8, 14.2} {
		lift := p.lift[i]
		limb{pt{x, 26}, pt{x, 35.6 - lift}, 3.6, denimBase}.draw(d)
		d.add(`<path d="M%.2f 27 L%.2f %.2f" stroke="%s" stroke-width="0.9"/>`, x+1.1, x+1.1, 35-lift, denimShade)
		d.add(`<rect x="%.2f" y="%.2f" width="4.2" height="1.3" rx="0.5" fill="%s" stroke="%s" stroke-width="0.5"/>`, x-2.1, 34.4-lift, denimShade, line)
		boot(d, pt{x, 40 - lift}, 0)
	}
	// Arms hanging at the sides, swinging, the right one holding the can.
	for i, x := range []float64{6.2, 17.8} {
		s := p.swing * []float64{1, -1}[i]
		hand := pt{x + []float64{-0.8, 0.8}[i], 27 + s}
		if p.can && i == 1 && !back || p.can && i == 0 && back {
			hand = pt{x + []float64{-3, 3}[i], 24.6}
		}
		elbow := pt{x + []float64{-0.5, 0.5}[i], 22.6}
		limb{pt{x, 19}, elbow, 2.8, shirtBase}.draw(d)
		limb{elbow, hand, 2.2, skin}.draw(d)
		d.add(`<path d="M%.2f %.2f l1.6 0.5" stroke="%s" stroke-width="0.9" transform="translate(-0.8 -0.3)"/>`, elbow.x, elbow.y, shirtShade)
		if p.can && hand.y < 25 {
			canAt(d, hand, map[bool]float64{true: 1, false: -1}[i == 1])
		}
	}
	// The shirt, and the overalls over it.
	shirt := d.linear(grad{0, 0, 1, 0, []stop{{0, shirtLight}, {0.6, shirtBase}, {1, shirtShade}}})
	d.add(`<path d="M6.6 18.6 Q12 16.6 17.4 18.6 L16.6 26.4 L7.4 26.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, shirt, line)
	denim := d.linear(grad{0, 0, 1, 0, []stop{{0, denimLight}, {0.45, denimBase}, {1, denimShade}}})
	d.add(`<path d="M7.6 24.6 L16.4 24.6 L16.8 28.6 L7.2 28.6 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	if back {
		// The straps cross between the shoulders.
		d.add(`<path d="M8.4 18.4 L15 24.8 M15.6 18.4 L9 24.8" stroke="%s" stroke-width="1.7"/>`, line)
		d.add(`<path d="M8.4 18.4 L15 24.8 M15.6 18.4 L9 24.8" stroke="%s" stroke-width="1"/>`, denimBase)
		for _, x := range []float64{8.6, 13.4} {
			d.add(`<rect x="%.2f" y="25.6" width="2" height="1.8" rx="0.4" fill="%s" stroke="%s" stroke-width="0.4"/>`, x, denimShade, line)
		}
	} else {
		d.add(`<path d="M8.8 20 L15.2 20 L15.6 25 L8.4 25 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
		d.add(`<rect x="10.2" y="21.2" width="3.6" height="2.4" rx="0.5" fill="%s" stroke="%s" stroke-width="0.4" stroke-dasharray="0.6 0.4"/>`, denimShade, denimLight)
		for _, s := range []float64{-1, 1} {
			d.add(`<path d="M%.2f 20.2 L%.2f 18.2" stroke="%s" stroke-width="1.6"/>`, 12+s*2.8, 12+s*3.6, line)
			d.add(`<path d="M%.2f 20.2 L%.2f 18.2" stroke="%s" stroke-width="0.9"/>`, 12+s*2.8, 12+s*3.6, denimBase)
			d.add(`<circle cx="%.2f" cy="20.4" r="0.65" fill="%s" stroke="%s" stroke-width="0.35"/>`, 12+s*2.8, brass, line)
		}
		d.add(`<path d="M10.4 17.6 L12 19.2 L13.6 17.6" stroke="%s" stroke-width="0.5" fill="%s"/>`, line, skin)
	}
	// The head: the face, or the back of the hair, and the hat over it.
	d.add(`<rect x="10.9" y="14.6" width="2.2" height="3" fill="%s" stroke="%s" stroke-width="0.5"/>`, skinShade, line)
	if back {
		d.add(`<ellipse cx="12" cy="10.6" rx="4.6" ry="5" fill="%s" stroke="%s" stroke-width="0.7"/>`, d.down(hairLight, hairShade), line)
		d.add(`<path d="M9.4 13.4 q1 1.8 2.6 2 q1.6 -0.2 2.6 -2 M10.4 9 q0.6 3 0.2 5 M13.6 9 q-0.6 3 -0.2 5" stroke="%s" stroke-width="0.5" fill="none"/>`, hairShade)
	} else {
		d.add(`<path d="M7.6 9.6 q-0.6 4 0.8 6 l1.2 -3.4 Z M16.4 9.6 q0.6 4 -0.8 6 l-1.2 -3.4 Z" fill="%s" stroke="%s" stroke-width="0.5"/>`, hairBase, line)
		faceAt(d, pt{12, 11})
		bangs(d, 7.7, 16.3, 7.6)
	}
	hat(d, 12, 7.2)
	return d.String()
}

// canAt is the watering can held at hand, its spout toward dir, pouring.
func canAt(d *doc, hand pt, dir float64) {
	x := hand.x + dir*1.6
	body := d.down(metal, metalShade)
	d.add(`<path d="M%.2f %.2f l%.2f -2.6" stroke="%s" stroke-width="1.5" stroke-linecap="round"/>`, x+dir*2, hand.y, dir*3, line)
	d.add(`<path d="M%.2f %.2f l%.2f -2.6" stroke="%s" stroke-width="0.7" stroke-linecap="round"/>`, x+dir*2, hand.y, dir*3, metal)
	d.add(`<rect x="%.2f" y="%.2f" width="4.6" height="3.8" rx="1" fill="%s" stroke="%s" stroke-width="0.6"/>`, x-2.3, hand.y-1.4, body, line)
	for i := range 3 {
		d.add(`<circle cx="%.2f" cy="%.2f" r="0.4" fill="%s"/>`, x+dir*(5.2+float64(i)*0.4), hand.y-2+float64(i)*1.6, waterLight)
	}
}
