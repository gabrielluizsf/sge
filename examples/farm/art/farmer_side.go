package art

// side is the farmer seen from their right, facing right: stride lifts
// the near leg forward, or back for a negative one, the far leg the
// other way.
func side(stride, swing float64, can bool) string {
	d := newDoc(FarmerW, FarmerH)
	d.add(`<ellipse cx="12" cy="39.2" rx="5" ry="1.3" fill="%s" opacity="0.3"/>`, shadowInk)
	// The far leg, shaded, then the near one.
	for i, s := range []float64{-stride, stride} {
		c := []string{denimShade, denimBase}[i]
		foot := pt{12 + s*3, 35.6 - abs(s)*0.6}
		limb{pt{12, 26}, foot, 3.6, c}.draw(d)
		boot(d, pt{foot.x, 40 - abs(s)*0.6}, 1)
	}
	// The body and the near arm.
	shirt := d.linear(grad{0, 0, 1, 0, []stop{{0, shirtShade}, {0.5, shirtBase}, {1, shirtLight}}})
	d.add(`<path d="M9.2 18.4 Q12 16.8 14.8 18.4 L14.6 26.4 L9.4 26.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, shirt, line)
	denim := d.linear(grad{0, 0, 1, 0, []stop{{0, denimShade}, {0.6, denimBase}, {1, denimLight}}})
	d.add(`<path d="M11.4 20 L14.8 20 L15 28.6 L9.2 28.6 L9.4 24.6 L11.4 24.6 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	d.add(`<path d="M12.6 20.2 L11 18.2" stroke="%s" stroke-width="1.5"/><path d="M12.6 20.2 L11 18.2" stroke="%s" stroke-width="0.8"/>`, line, denimBase)
	d.add(`<circle cx="12.8" cy="20.5" r="0.65" fill="%s" stroke="%s" stroke-width="0.35"/>`, brass, line)
	hand := pt{12 + swing*2.4, 27}
	if can {
		hand = pt{16.4, 24.6}
	}
	elbow := pt{12 + swing, 22.8}
	limb{pt{11.8, 19}, elbow, 2.8, shirtBase}.draw(d)
	limb{elbow, hand, 2.2, skin}.draw(d)
	if can {
		canAt(d, hand, 1)
	}
	// The head in profile: the back of the hair, the face to the right.
	d.add(`<rect x="10.8" y="14.6" width="2.2" height="3" fill="%s" stroke="%s" stroke-width="0.5"/>`, skinShade, line)
	d.add(`<path d="M15.6 9 Q17.2 11.4 16.4 13.6 Q15.4 15.6 12.6 15.6 Q10.2 15.4 9.6 13 L10 9 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, skin, line)
	d.add(`<path d="M16.3 11.6 l0.9 0.5 l-0.8 0.3" stroke="%s" stroke-width="0.45" fill="%s"/>`, line, skin)
	eyeAt(d, pt{14.6, 11.4})
	d.add(`<ellipse cx="14" cy="13.4" rx="0.9" ry="0.45" fill="%s" opacity="0.55"/>`, blush)
	d.add(`<path d="M15.2 14.4 q0.5 0.3 1 0" stroke="%s" stroke-width="0.4" fill="none"/>`, line)
	d.add(`<path d="M13.6 7.6 Q8.4 6.4 7.8 10.6 Q7.6 13.6 9.4 15.6 L10.8 13 L11.6 14.4 L12.2 11.4 Q13.6 10.6 14.8 9.6 L15.8 9.8 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, d.down(hairLight, hairBase), line)
	d.add(`<ellipse cx="11.2" cy="12" rx="0.9" ry="1.2" fill="%s" stroke="%s" stroke-width="0.45"/>`, skin, line)
	hat(d, 12.6, 7.2)
	return d.String()
}
