package flock

import (
	"math"
)

// takeOff has b rise toward another perch, a good way off, if the
// perches near it have one; otherwise it stays.
func (f *Flock) takeOff(b *Bird) {
	ps := f.Perches(b.Pos, f.Range/2)
	for range 8 {
		if len(ps) == 0 {
			break
		}
		p := ps[b.src.Intn(len(ps))]
		if d := p.At.Distance(b.Pos); d < f.Range/8 || f.taken(b, p) {
			continue
		}
		b.from = Perch{At: b.Pos, Height: b.Alt}
		b.to, b.State, b.Time, b.leg = p, Flying, 0, 0
		b.span = p.At.Distance(b.Pos)
		if d := p.At.Sub(b.Pos); d.Length() > 0 {
			b.Dir = d.Scale(1 / d.Length())
		}
		return
	}
	b.rest = 3 + 6*b.src.Float() // nowhere to go: it sits a while longer
}

// fly moves b dt seconds along its flight, climbing and settling as it
// goes, and lands it on its perch at the end.
func (f *Flock) fly(b *Bird, dt float64) {
	b.leg = math.Min(b.leg+f.Speed*dt/math.Max(b.span, 1), 1)
	b.Pos = b.from.At.Lerp(b.to.At, b.leg)
	arc := math.Sin(math.Pi * b.leg)
	b.Alt = b.from.Height + (b.to.Height-b.from.Height)*b.leg + f.Lift*arc
	if b.leg >= 1 {
		b.Pos, b.Alt, b.State, b.Time = b.to.At, b.to.Height, Perched, 0
		b.rest = 8 + 16*b.src.Float()
	}
}
