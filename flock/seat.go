package flock

import "github.com/gabrielluizsf/sge/vector"

// seat puts b on a perch within Range of c that no other bird sits near,
// if there is one; it reports whether it did.
func (f *Flock) seat(b *Bird, c vector.Vector3) bool {
	ps := f.Perches(c, f.Range)
	for range 8 {
		if len(ps) == 0 {
			return false
		}
		p := ps[b.src.Intn(len(ps))]
		if f.taken(b, p) {
			continue
		}
		b.Pos, b.Alt, b.State, b.Time = p.At, p.Height, Perched, 0
		b.from, b.to = p, p
		b.rest = 6 + 14*b.src.Float()
		return true
	}
	return false
}

// taken reports whether a bird other than b sits on, or flies to, p's
// neighbourhood.
func (f *Flock) taken(b *Bird, p Perch) bool {
	for _, o := range f.birds {
		if o != b && o.to.At.Distance(p.At) < 48 {
			return true
		}
	}
	return false
}
