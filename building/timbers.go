package building

// timberW is how wide a timber of a frame shows on the wall.
const timberW = 0.18

// timbers are the frame showing on side s of a timbered house: a post at
// each corner, a beam along its foot, each floor and its top, a stud
// either side of each door and window and between the bays, and braces
// slanting across the panels at its ends.
func (r *raising) timbers(s side) {
	b, half := r.b, s.length/2
	wood := paint{of: Timber}
	rail := func(y float64) {
		b.block(s.box(patch{lo: vec{-half, y - timberW/2, 0}, hi: vec{half, y + timberW/2, 0.045}, paint: wood}))
	}
	rail(floor + timberW/2)
	for k := 1; k < r.h.Storeys; k++ {
		rail(floor + float64(k)*r.h.Storey)
	}
	rail(r.eaves - timberW/2)
	stud := func(x float64, storey int) {
		lo := floor + float64(storey)*r.h.Storey
		b.block(s.box(patch{lo: vec{x - timberW/2, lo, 0}, hi: vec{x + timberW/2, lo + r.h.Storey, 0.035}, paint: wood}))
	}
	at := bays(s)
	for storey := range r.h.Storeys {
		stud(-half+timberW/2, storey)
		stud(half-timberW/2, storey)
		for i, x := range at {
			w := winW/2 + 0.08 + timberW/2
			if storey == 0 && s.front() && i == doorBay(len(at)) {
				w = doorW/2 + 0.1 + timberW/2
			}
			stud(x-w, storey)
			stud(x+w, storey)
			if i > 0 {
				stud((x+at[i-1])/2, storey)
			}
		}
		r.braces(s, storey)
	}
}

// braces slant across the panels at side s's ends, in storey storey,
// from a corner post's foot up to the stud by the first window in.
func (r *raising) braces(s side, storey int) {
	at := bays(s)
	half := s.length / 2
	lo := floor + float64(storey)*r.h.Storey + timberW
	hi := floor + float64(storey+1)*r.h.Storey - timberW
	reach := at[0] - (winW/2 + 0.08 + timberW) - (-half + timberW)
	if reach < 0.4 {
		return
	}
	for _, dir := range []float64{-1, 1} {
		if dir < 0 && storey == 0 && s.front() && doorBay(len(at)) == 0 {
			continue // the door is in the way
		}
		from := s.at(vec{dir * (half - timberW), lo, 0.02})
		to := s.at(vec{dir * (half - timberW - reach), hi, 0.02})
		r.b.lay(beam{a: from, b: to, w: 0.04, h: timberW * 0.9, paint: paint{of: Timber}})
	}
}
