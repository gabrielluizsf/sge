package flock

import "github.com/gabrielluizsf/sge/vector"

// Step lets the flock live dt seconds: birds are made and kept about the
// view, perched ones rise when someone is near, flying ones go on.
func (f *Flock) Step(dt float64) {
	if f.Around == nil || f.Perches == nil {
		return
	}
	around := f.Around()
	for len(f.birds) < f.Count {
		b := &Bird{src: f.source(), Dir: vector.Vector3{X: 1}}
		f.birds = append(f.birds, b)
		if !f.seat(b, around) {
			f.birds = f.birds[:len(f.birds)-1]
			return
		}
	}
	var scares []vector.Vector3
	if f.Scares != nil {
		scares = f.Scares()
	}
	for _, b := range f.birds {
		b.Time += dt
		switch {
		case b.State == Flying:
			f.fly(b, dt)
		case near(b.Pos, scares, f.Fear):
			f.takeOff(b)
		case b.Pos.Distance(around) > 1.5*f.Range:
			f.seat(b, around)
		default:
			if b.rest -= dt; b.rest < 0 {
				f.takeOff(b)
			}
		}
	}
}

// near reports whether any of the points is within d of p.
func near(p vector.Vector3, points []vector.Vector3, d float64) bool {
	for _, q := range points {
		if p.Distance(q) < d {
			return true
		}
	}
	return false
}
