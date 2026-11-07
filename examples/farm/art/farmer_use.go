package art

// What a person does with their hands besides swinging, the moments
// drawn: eating, with the food at the chest or at the mouth; filling the
// can at a well, the free hand reaching up the rope, pulling it down or
// tipping the pail over the can; and at a river's bank, crouched, the
// can sunk in the water, or lifted out of it.
const (
	eatLow = 1 + iota
	eatUp
	wellReach
	wellPull
	wellPour
	dipDown
	dipUp
	milkHigh
	milkLow
)

// crouched is whether the pose is on the heels, at the water's edge.
func (p pose) crouched() bool {
	return p.use == dipDown || p.use == dipUp || p.use == milkHigh || p.use == milkLow
}

// squat starts the group that draws a crouched body: the whole drawing
// pressed down onto the ground, which is where the feet still stand.
func (p pose) squat(d *doc) {
	if p.crouched() {
		d.add(`<g transform="translate(0 9.6) scale(1 0.76)">`)
	}
}

// stand ends what squat started.
func (p pose) stand(d *doc) {
	if p.crouched() {
		d.add(`</g>`)
	}
}

// useHand is where an arm's hand is in the pose's use, seen from the
// front or the back, the arm at x and out to dir, -1 left or 1 right of
// the picture, and whether it holds what is used, the food or the can.
// ok is false when the pose does not use the arm.
func (p pose) useHand(x, dir float64, holds bool) (hand pt, ok bool) {
	in := x - dir*3.4
	switch {
	case p.use == eatLow && holds:
		return pt{in, 22.4}, true
	case p.use == eatUp && holds:
		return pt{in, 14.6}, true
	case p.use == dipDown && holds:
		return pt{in, 31}, true
	case p.use == dipUp && holds:
		return pt{in, 27.4}, true
	case p.use == dipDown || p.use == dipUp:
		return pt{x - dir*0.4, 29}, true // the free hand flat on the bank
	case p.use == milkHigh && holds, p.use == milkLow && !holds:
		return pt{in, 26}, true
	case p.use == milkHigh, p.use == milkLow:
		return pt{in, 28.6}, true
	case p.use == wellReach:
		return pt{x + dir*0.6, 15.8}, true
	case p.use == wellPull:
		return pt{x - dir*0.4, 22.4}, true
	case p.use == wellPour:
		return pt{x - dir*1.4, 18.8}, true
	}
	return pt{}, false
}

// useProp draws what the hand holds in the pose's use, held at hand: the
// food, or the can, with the water that goes into it. dir is the way the
// can's spout points, -1 left or 1 right.
func (p pose) useProp(d *doc, hand pt, dir float64) {
	switch p.use {
	case eatLow, eatUp:
		apple(d, hand)
	case wellPour:
		waterFall(d, pt{hand.x + dir*1.6, hand.y - 1.4})
	case dipDown:
		d.add(`<ellipse cx="%.2f" cy="%.2f" rx="4.6" ry="1.2" fill="%s" opacity="0.7"/>`, hand.x+dir*1.6, hand.y+1.4, waterBase)
		canAt(d, hand, dir)
	case dipUp:
		canAt(d, hand, dir)
	case milkHigh, milkLow:
		pail(d, pt{hand.x, 33})
	}
}

// pail is the milk pail on the ground at at, full of white.
func pail(d *doc, at pt) {
	d.add(`<path d="M%.2f %.2f l5.2 0 l-0.8 5 l-3.6 0 Z" fill="%s" stroke="%s" stroke-width="0.6" stroke-linejoin="round"/>`, at.x-2.6, at.y, metal, line)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="2.6" ry="0.8" fill="%s" stroke="%s" stroke-width="0.4"/>`, at.x, at.y, white, line)
}

// apple is a bite of food in the hand at at: a small red apple with a
// stem.
func apple(d *doc, at pt) {
	d.add(`<circle cx="%.2f" cy="%.2f" r="1.5" fill="%s" stroke="%s" stroke-width="0.5"/>`, at.x, at.y-0.6, "#c8372c", line)
	d.add(`<path d="M%.2f %.2f l0.3 -0.9" stroke="%s" stroke-width="0.4"/>`, at.x, at.y-1.9, line)
	d.add(`<circle cx="%.2f" cy="%.2f" r="0.4" fill="%s" opacity="0.8"/>`, at.x-0.5, at.y-1, white)
}

// waterFall is water running down into the can's mouth at at: a short
// stream of drops from above.
func waterFall(d *doc, at pt) {
	d.add(`<path d="M%.2f %.2f L%.2f %.2f" stroke="%s" stroke-width="1" stroke-linecap="round"/>`, at.x, at.y-6, at.x, at.y, waterLight)
	d.add(`<path d="M%.2f %.2f L%.2f %.2f" stroke="%s" stroke-width="0.4" stroke-linecap="round" opacity="0.9"/>`, at.x-0.2, at.y-6, at.x-0.2, at.y, white)
}

// sideHand is where the near arm's hand is in the pose's use, seen from
// the side, the person facing right: the arm that holds the food or the
// can, the only one seen.
func (p pose) sideHand() (hand pt, ok bool) {
	switch p.use {
	case eatLow:
		return pt{14.6, 22.2}, true
	case eatUp:
		return pt{16.4, 14.4}, true
	case dipDown:
		return pt{17.6, 31}, true
	case dipUp:
		return pt{16.4, 27.4}, true
	case milkHigh:
		return pt{18.4, 26}, true
	case milkLow:
		return pt{18.4, 28.6}, true
	}
	return pt{}, false
}

// foreArm draws, over the body, the arm that brings the food up, seen
// from the front or the back: the arm is in front of the chest, where
// the first drawing of it is under the shirt.
func (p pose) foreArm(d *doc, l Look, back bool) {
	if p.use != eatLow && p.use != eatUp {
		return
	}
	x, dir := 17.8, 1.0
	if back {
		x, dir = 6.2, -1.0
	}
	hand, _ := p.useHand(x, dir, true)
	elbow := pt{x - dir*1.2, 24}
	if p.use == eatUp {
		elbow = pt{x - dir*0.4, 23}
	}
	limb{pt{x, 19}, elbow, 2.8, l.Shirt[1]}.draw(d)
	limb{elbow, hand, 2.2, l.Skin[0]}.draw(d)
	apple(d, hand)
}
