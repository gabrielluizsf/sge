package particle

import "github.com/gabrielluizsf/sge/vector"

// Emit throws n specks from the point at, at height z, each with the
// system's own spread added to push, a velocity over the ground.
func (s *System) Emit(at vector.Vector3, push vector.Vector3, n int) {
	if len(s.specks) == 0 {
		return
	}
	c := &s.cfg
	for range n {
		p := &s.specks[s.next]
		if s.next++; s.next == len(s.specks) {
			s.next = 0
		}
		if !p.live {
			s.live++
		}
		p.live = true
		p.x, p.y, p.z = at.X, at.Y, at.Z
		p.vx = push.X + (s.rng.Float()*2-1)*c.Speed
		p.vy = push.Y + (s.rng.Float()*2-1)*c.Speed
		p.vz = c.Lift * (0.5 + s.rng.Float()*0.5)
		p.lain = 0
		p.rest = c.Rest + (s.rng.Float()*2-1)*c.Spread
		p.size = c.Size + s.rng.Intn(c.Grow+1)
		p.color = c.Colors[s.rng.Intn(len(c.Colors))]
	}
}
