package building

// side is a face of a wall, as things are put on it: where its middle
// meets the ground, the way along it to the right as one faces it, the
// way out of it, and how long it is.
type side struct {
	mid, right, out vec
	length          float64
}

// sides are the four faces of walls w wide along X and d deep along Z,
// standing about the origin: the front, the right, the back, the left.
func sides(w, d float64) [4]side {
	return [4]side{
		{mid: vec{0, 0, d / 2}, right: vec{1, 0, 0}, out: vec{0, 0, 1}, length: w},
		{mid: vec{w / 2, 0, 0}, right: vec{0, 0, -1}, out: vec{1, 0, 0}, length: d},
		{mid: vec{0, 0, -d / 2}, right: vec{-1, 0, 0}, out: vec{0, 0, -1}, length: w},
		{mid: vec{-w / 2, 0, 0}, right: vec{0, 0, 1}, out: vec{-1, 0, 0}, length: d},
	}
}

// front is whether the side faces the front, +Z, where the door is.
func (s side) front() bool { return s.out[2] > 0.5 }

// at is the point l of the side: l[0] along it from its middle, l[1] up
// and l[2] out of it.
func (s side) at(l vec) vec {
	return add(add(s.mid, scale(s.right, l[0])), add(scale(up, l[1]), scale(s.out, l[2])))
}

// patch is a box on the side, from lo to hi as at takes them, of paint p.
type patch struct {
	lo, hi vec
	paint  paint
}

// box is the patch as a box in the building's space; wood's grain runs
// along it, level if it is long and low, else up, as boards and beams
// run; other textures stand as they are, their down down.
func (s side) box(p patch) box {
	a, b := s.at(p.lo), s.at(p.hi)
	lo := vec{min(a[0], b[0]), min(a[1], b[1]), min(a[2], b[2])}
	hi := vec{max(a[0], b[0]), max(a[1], b[1]), max(a[2], b[2])}
	var grain vec
	if p.paint.of == Timber || p.paint.of == Planks {
		grain = up
		if p.hi[0]-p.lo[0] > p.hi[1]-p.lo[1] {
			grain = s.right
		}
	}
	return box{lo: lo, hi: hi, paint: p.paint, grain: grain}
}
