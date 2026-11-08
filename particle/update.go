package particle

import "github.com/gabrielluizsf/sge/game"

// Update moves every speck dt seconds on: those in the air fly and fall,
// and on landing stop where they are; those on the ground lie until their
// rest is over, and are gone.
func (s *System) Update(_ *game.Object, dt float64) {
	g := s.cfg.Gravity
	for i := range s.specks {
		p := &s.specks[i]
		if !p.live {
			continue
		}
		if p.z > 0 || p.vz > 0 {
			p.vz -= g * dt
			p.x += p.vx * dt
			p.y += p.vy * dt
			if p.z += p.vz * dt; p.z <= 0 {
				p.z, p.vx, p.vy, p.vz = 0, 0, 0, 0
			}
			continue
		}
		if p.lain += dt; p.lain >= p.rest {
			p.live = false
			s.live--
		}
	}
}
