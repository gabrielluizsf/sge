package building

import "image/color"

// Doors and windows, in metres: a door is 0.95 wide and 2.1 tall, a
// window 0.9 by 1.2, its sill 0.9 above its floor; a bay, the stretch of
// wall each takes, is at least bay wide.
const (
	doorW, doorH = 0.95, 2.1
	winW, winH   = 0.9, 1.2
	sillAt       = 0.9
	bay          = 2.4
)

// bays are the middles of the bays along side s, from its left.
func bays(s side) []float64 {
	n := max(1, int(s.length/bay))
	out := make([]float64, n)
	for i := range out {
		out[i] = -s.length/2 + (float64(i)+0.5)*s.length/float64(n)
	}
	return out
}

// doorBay is the bay of the front door: the middle one, or, of an even
// number, the one left of the middle.
func doorBay(n int) int { return (n - 1) / 2 }

// openings are side s's windows, a row a storey, and on the front its
// door, in the ground floor's door bay, with steps up to it.
func (r *raising) openings(s side, front bool) {
	at := bays(s)
	for storey := range r.h.Storeys {
		base := floor + float64(storey)*r.h.Storey
		for i, x := range at {
			if front && storey == 0 && i == doorBay(len(at)) {
				r.door(s, x)
				continue
			}
			r.window(s, vec{x, base + sillAt})
		}
	}
}

// trim is the paint of frames: painted, or bare timber.
func (r *raising) trim() paint { return coat(Timber, r.h.Trim) }

// coat is material m painted c, or bare if c is none.
func coat(m Material, c color.NRGBA) paint {
	return paint{of: m, tint: c, painted: c.A != 0}
}

// door is a door on side s, x along it: a plank door in a frame, its
// iron handle, and, without a porch, stone steps up to it.
func (r *raising) door(s side, x float64) {
	b, h := r.b, doorH
	lo, hi := x-doorW/2, x+doorW/2
	b.block(s.box(patch{lo: vec{lo - 0.1, floor, 0}, hi: vec{lo, floor + h + 0.05, 0.08}, paint: r.trim()}))
	b.block(s.box(patch{lo: vec{hi, floor, 0}, hi: vec{hi + 0.1, floor + h + 0.05, 0.08}, paint: r.trim()}))
	b.block(s.box(patch{lo: vec{lo - 0.16, floor + h, 0}, hi: vec{hi + 0.16, floor + h + 0.16, 0.09}, paint: r.trim()}))
	b.block(s.box(patch{lo: vec{lo, floor, 0}, hi: vec{hi, floor + h, 0.035}, paint: coat(Planks, r.h.Door)}))
	b.block(s.box(patch{lo: vec{hi - 0.16, floor + 1.0, 0.035}, hi: vec{hi - 0.1, floor + 1.06, 0.075}, paint: paint{of: Iron}}))
	if r.h.Porch {
		return
	}
	stone := paint{of: Stone}
	b.block(s.box(patch{lo: vec{lo - 0.3, 0, 0}, hi: vec{hi + 0.3, floor * 0.5, 0.65}, paint: stone}))
	b.block(s.box(patch{lo: vec{lo - 0.3, 0, 0}, hi: vec{hi + 0.3, floor, 0.33}, paint: stone}))
}

// window is a window on side s, its sill's middle at at: its glass, its
// frame and the cross dividing its panes, a stone sill under it, and
// shutters folded back beside it, if the house has them.
func (r *raising) window(s side, at vec) {
	b, x, y := r.b, at[0], at[1]
	lo, hi := x-winW/2, x+winW/2
	b.block(s.box(patch{lo: vec{lo, y, 0}, hi: vec{hi, y + winH, 0.02}, paint: paint{of: Glass}}))
	f := 0.08
	for _, p := range []patch{
		{lo: vec{lo - f, y, 0}, hi: vec{lo, y + winH, 0.06}}, {lo: vec{hi, y, 0}, hi: vec{hi + f, y + winH, 0.06}},
		{lo: vec{lo - f, y + winH, 0}, hi: vec{hi + f, y + winH + f, 0.07}}, {lo: vec{lo - f, y - 0.02, 0}, hi: vec{hi + f, y, 0.06}},
		{lo: vec{x - 0.025, y, 0}, hi: vec{x + 0.025, y + winH, 0.045}}, {lo: vec{lo, y + winH*0.62, 0}, hi: vec{hi, y + winH*0.62 + 0.05, 0.05}},
	} {
		p.paint = r.trim()
		b.block(s.box(p))
	}
	b.block(s.box(patch{lo: vec{lo - 0.14, y - 0.09, 0}, hi: vec{hi + 0.14, y - 0.02, 0.13}, paint: paint{of: Stone}}))
	if r.h.Shutters == nil {
		return
	}
	shutter := coat(Planks, *r.h.Shutters)
	w := winW/2 + 0.02
	for _, edge := range []float64{lo - f - w, hi + f} {
		b.block(s.box(patch{lo: vec{edge, y - 0.02, 0}, hi: vec{edge + w, y + winH + 0.04, 0.045}, paint: shutter}))
	}
}
