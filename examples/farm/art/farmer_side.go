package art

// side is a person of look l seen from their right, facing right: the
// pose's stride lifts the near leg forward, or back for a negative one,
// the far leg the other way.
func side(l Look, p pose) string {
	stride, swing := p.stride, p.swing
	d := newDoc(FarmerW, FarmerH)
	d.add(`<ellipse cx="12" cy="39.2" rx="5" ry="1.3" fill="%s" opacity="0.3"/>`, shadowInk)
	headGroup(d)
	hairBehind(d, l, true)
	endGroup(d)
	bodyGroup(d)
	// The far leg, shaded, then the near one.
	for i, s := range []float64{-stride, stride} {
		foot := pt{12 + s*3, 35.6 - abs(s)*0.6}
		if l.Skirt {
			limb{pt{12, 28}, foot, 2.6, []string{l.Skin[1], l.Skin[0]}[i]}.draw(d)
		} else {
			limb{pt{12, 26}, foot, 3.6, []string{l.Pants[2], l.Pants[1]}[i]}.draw(d)
		}
		boot(d, pt{foot.x, 40 - abs(s)*0.6}, 1)
	}
	// The body and the near arm.
	shirt := d.linear(grad{0, 0, 1, 0, []stop{{0, l.Shirt[2]}, {0.5, l.Shirt[1]}, {1, l.Shirt[0]}}})
	d.add(`<path d="M9.2 18.4 Q12 16.8 14.8 18.4 L14.6 26.4 L9.4 26.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, shirt, line)
	denim := d.linear(grad{0, 0, 1, 0, []stop{{0, l.Pants[2]}, {0.6, l.Pants[1]}, {1, l.Pants[0]}}})
	switch {
	case l.Overalls:
		d.add(`<path d="M11.4 20 L14.8 20 L15 28.6 L9.2 28.6 L9.4 24.6 L11.4 24.6 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
		d.add(`<path d="M12.6 20.2 L11 18.2" stroke="%s" stroke-width="1.5"/><path d="M12.6 20.2 L11 18.2" stroke="%s" stroke-width="0.8"/>`, line, l.Pants[1])
		d.add(`<circle cx="12.8" cy="20.5" r="0.65" fill="%s" stroke="%s" stroke-width="0.35"/>`, brass, line)
	case l.Skirt:
		d.add(`<path d="M9.4 24.6 L14.6 24.6 L16.4 31.4 L7.8 31.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	default:
		d.add(`<path d="M9.4 24.6 L14.8 24.6 L15 28.6 L9.2 28.6 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	}
	hand := pt{12 + swing*2.4, 27}
	if p.can {
		hand = pt{16.4, 24.6}
	}
	switch p.axe {
	case windUp:
		hand = pt{10, 15.4}
	case strike:
		hand = pt{17.6, 24}
	}
	elbow := pt{12 + swing, 22.8}
	limb{pt{11.8, 19}, elbow, 2.8, l.Shirt[1]}.draw(d)
	limb{elbow, hand, 2.2, l.Skin[0]}.draw(d)
	if p.can {
		canAt(d, hand, 1)
	}
	switch p.axe {
	case windUp:
		axeAt(d, hand, -60)
	case strike:
		swoosh(d, [3]pt{{8, 8}, {20, 6}, {21.6, 30}})
		axeAt(d, hand, 125)
	}
	endGroup(d)
	// The head in profile: the back of the hair, the face to the right.
	headGroup(d)
	d.add(`<rect x="10.8" y="14.6" width="2.2" height="3" fill="%s" stroke="%s" stroke-width="0.5"/>`, l.Skin[1], line)
	d.add(`<path d="M15.6 9 Q17.2 11.4 16.4 13.6 Q15.4 15.6 12.6 15.6 Q10.2 15.4 9.6 13 L10 9 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, l.Skin[0], line)
	d.add(`<path d="M16.3 11.6 l0.9 0.5 l-0.8 0.3" stroke="%s" stroke-width="0.45" fill="%s"/>`, line, l.Skin[0])
	eyeAt(d, pt{14.6, 11.4})
	d.add(`<ellipse cx="14" cy="13.4" rx="0.9" ry="0.45" fill="%s" opacity="0.55"/>`, blush)
	d.add(`<path d="M15.2 14.4 q0.5 0.3 1 0" stroke="%s" stroke-width="0.4" fill="none"/>`, line)
	hairSide(d, l)
	d.add(`<ellipse cx="11.2" cy="12" rx="0.9" ry="1.2" fill="%s" stroke="%s" stroke-width="0.45"/>`, l.Skin[0], line)
	if l.Hat {
		hat(d, 12.6, 7.2)
	}
	endGroup(d)
	return d.String()
}
