package art

import "fmt"

// pose is how a person stands in a frame: how far each leg is lifted
// mid-stride, seen from the front, or how far the near leg strides, seen
// from the side; how far the arms swing; and whether they hold what they
// use out, the can or the axe raised.
type pose struct {
	lift   [2]float64 // left leg, right leg, as the viewer sees them
	stride float64
	swing  float64
	can    bool
	axe    int // 0 no axe, windUp raised, strike swung down
}

// The swing of an axe: raised behind the head, then brought down.
const (
	windUp = 1 + iota
	strike
)

// front is a person of look l facing the viewer, or with back true facing
// away.
func front(l Look, p pose, back bool) string {
	d := newDoc(FarmerW, FarmerH)
	d.add(`<ellipse cx="12" cy="39.2" rx="6" ry="1.3" fill="%s" opacity="0.3"/>`, shadowInk)
	if !back {
		headGroup(d)
		hairBehind(d, l, false)
		endGroup(d)
	}
	bodyGroup(d)
	legs(d, l, p)
	// Arms hanging at the sides, swinging, the right one holding the can.
	for i, x := range []float64{6.2, 17.8} {
		s := p.swing * []float64{1, -1}[i]
		hand := pt{x + []float64{-0.8, 0.8}[i], 27 + s}
		holds := i == 1 && !back || i == 0 && back
		if p.can && holds {
			hand = pt{x + []float64{-3, 3}[i], 24.6}
		}
		if p.axe != 0 && holds {
			hand = swingHand(p.axe, back)
		}
		elbow := pt{x + []float64{-0.5, 0.5}[i], 22.6}
		limb{pt{x, 19}, elbow, 2.8, l.Shirt[1]}.draw(d)
		limb{elbow, hand, 2.2, l.Skin[0]}.draw(d)
		d.add(`<path d="M%.2f %.2f l1.6 0.5" stroke="%s" stroke-width="0.9" transform="translate(-0.8 -0.3)"/>`, elbow.x, elbow.y, l.Shirt[2])
		if p.can && hand.y < 25 {
			canAt(d, hand, map[bool]float64{true: 1, false: -1}[i == 1])
		}
		if p.axe != 0 && holds {
			swingAxe(d, p.axe, hand, back)
		}
	}
	// The shirt, and the overalls or trousers over it.
	shirt := d.linear(grad{0, 0, 1, 0, []stop{{0, l.Shirt[0]}, {0.6, l.Shirt[1]}, {1, l.Shirt[2]}}})
	d.add(`<path d="M6.6 18.6 Q12 16.6 17.4 18.6 L16.6 26.4 L7.4 26.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, shirt, line)
	denim := d.linear(grad{0, 0, 1, 0, []stop{{0, l.Pants[0]}, {0.45, l.Pants[1]}, {1, l.Pants[2]}}})
	if l.Skirt {
		d.add(`<path d="M7.6 24.6 L16.4 24.6 L18 31.4 L6 31.4 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
		d.add(`<path d="M10 25.4 L9.4 31 M14 25.4 L14.6 31" stroke="%s" stroke-width="0.5"/>`, l.Pants[2])
	} else {
		d.add(`<path d="M7.6 24.6 L16.4 24.6 L16.8 28.6 L7.2 28.6 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	}
	if l.Overalls {
		overalls(d, l, denim, back)
	} else if !back {
		d.add(`<path d="M10.4 17.6 L12 19.2 L13.6 17.6" stroke="%s" stroke-width="0.5" fill="%s"/>`, line, l.Skin[0])
	}
	endGroup(d)
	// The head: the face, or the back of the hair, and the hat over it.
	headGroup(d)
	d.add(`<rect x="10.9" y="14.6" width="2.2" height="3" fill="%s" stroke="%s" stroke-width="0.5"/>`, l.Skin[1], line)
	if back {
		hairBack(d, l)
	} else {
		hairSides(d, l)
		faceAt(d, pt{12, 11}, l)
		hairTop(d, l)
	}
	if l.Hat {
		hat(d, 12, 7.2)
	}
	endGroup(d)
	return d.String()
}

// legs are a person's legs and boots seen from the front or back, each
// lifted as the stride says: in trousers, or bare under a skirt.
func legs(d *doc, l Look, p pose) {
	for i, x := range []float64{9.8, 14.2} {
		lift := p.lift[i]
		if l.Skirt {
			limb{pt{x, 28}, pt{x, 36 - lift}, 2.6, l.Skin[0]}.draw(d)
		} else {
			limb{pt{x, 26}, pt{x, 35.6 - lift}, 3.6, l.Pants[1]}.draw(d)
			d.add(`<path d="M%.2f 27 L%.2f %.2f" stroke="%s" stroke-width="0.9"/>`, x+1.1, x+1.1, 35-lift, l.Pants[2])
			d.add(`<rect x="%.2f" y="%.2f" width="4.2" height="1.3" rx="0.5" fill="%s" stroke="%s" stroke-width="0.5"/>`, x-2.1, 34.4-lift, l.Pants[2], line)
		}
		boot(d, pt{x, 40 - lift}, 0)
	}
}

// overalls are the bib and straps of overalls, in the paint denim, their
// straps crossed on the back.
func overalls(d *doc, l Look, denim string, back bool) {
	if back {
		d.add(`<path d="M8.4 18.4 L15 24.8 M15.6 18.4 L9 24.8" stroke="%s" stroke-width="1.7"/>`, line)
		d.add(`<path d="M8.4 18.4 L15 24.8 M15.6 18.4 L9 24.8" stroke="%s" stroke-width="1"/>`, l.Pants[1])
		for _, x := range []float64{8.6, 13.4} {
			d.add(`<rect x="%.2f" y="25.6" width="2" height="1.8" rx="0.4" fill="%s" stroke="%s" stroke-width="0.4"/>`, x, l.Pants[2], line)
		}
		return
	}
	d.add(`<path d="M8.8 20 L15.2 20 L15.6 25 L8.4 25 Z" fill="%s" stroke="%s" stroke-width="0.7" stroke-linejoin="round"/>`, denim, line)
	d.add(`<rect x="10.2" y="21.2" width="3.6" height="2.4" rx="0.5" fill="%s" stroke="%s" stroke-width="0.4" stroke-dasharray="0.6 0.4"/>`, l.Pants[2], l.Pants[0])
	for _, s := range []float64{-1, 1} {
		d.add(`<path d="M%.2f 20.2 L%.2f 18.2" stroke="%s" stroke-width="1.6"/>`, 12+s*2.8, 12+s*3.6, line)
		d.add(`<path d="M%.2f 20.2 L%.2f 18.2" stroke="%s" stroke-width="0.9"/>`, 12+s*2.8, 12+s*3.6, l.Pants[1])
		d.add(`<circle cx="%.2f" cy="20.4" r="0.65" fill="%s" stroke="%s" stroke-width="0.35"/>`, 12+s*2.8, brass, line)
	}
	d.add(`<path d="M10.4 17.6 L12 19.2 L13.6 17.6" stroke="%s" stroke-width="0.5" fill="%s"/>`, line, l.Skin[0])
}

// swingHand is where the hand holding the axe is at the swing's moment,
// seen from the front, or from the back.
func swingHand(at int, back bool) pt {
	switch {
	case back && at == windUp:
		return pt{4.6, 13}
	case back:
		return pt{9.6, 21.4}
	case at == windUp:
		return pt{19.4, 13}
	}
	return pt{15, 27.4}
}

// swingAxe is the axe in hand at the swing's moment, seen from the front,
// or the back, where the swing down goes away from the viewer, hidden.
func swingAxe(d *doc, at int, hand pt, back bool) {
	switch {
	case back && at == windUp:
		axeAt(d, hand, 35)
	case at == windUp:
		axeAt(d, hand, -35)
	case !back:
		swoosh(d, [3]pt{{21, 9}, {23, 22}, {16.4, 34}})
		axeAt(d, hand, 165)
	}
}

// axeAt is the axe held at hand, turned deg degrees from pointing up, its
// blade to the right of its handle as it points.
func axeAt(d *doc, hand pt, deg float64) {
	d.add(`<g transform="translate(%.2f %.2f) rotate(%g)">`, hand.x, hand.y, deg)
	limb{pt{0, 1.6}, pt{0, -7.4}, 1.1, woodBase}.draw(d)
	d.add(`<path d="M-0.6 -8.6 q2.8 -1.6 4.2 1.2 q-0.6 2 -4.2 2.8 Z" fill="%s" stroke="%s" stroke-width="0.5" stroke-linejoin="round"/>`, d.down(metal, metalShade), line)
	d.add(`<path d="M1.4 -8.6 q1.4 0.4 2 1.4" stroke="%s" stroke-width="0.4" fill="none"/>`, white)
	d.add(`</g>`)
}

// swoosh is the trail of a swing, through its three points: a white arc,
// thick in its middle, fading.
func swoosh(d *doc, p [3]pt) {
	path := fmt.Sprintf("M%.2f %.2f Q%.2f %.2f %.2f %.2f", p[0].x, p[0].y, p[1].x, p[1].y, p[2].x, p[2].y)
	d.add(`<path d="%s" stroke="%s" stroke-width="2.4" fill="none" stroke-linecap="round" opacity="0.35"/>`, path, white)
	d.add(`<path d="%s" stroke="%s" stroke-width="0.9" fill="none" stroke-linecap="round" opacity="0.8"/>`, path, white)
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
