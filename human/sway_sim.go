package human

import "math"

// The hair's spring: its own rate, in radians a second, how damped it is,
// how far, in radians, it may be swung, how long the hair is for the push
// of the head's motion on it, in metres, and how strong the breeze is, in
// radians.
const (
	swayRate    = 6.0
	swayDamping = 0.35
	swayLimit   = 0.14
	swayHang    = 0.25
	swayBreeze  = 0.035
	swayStep    = 1.0 / 120
)

// hairSwing is the hair's free ends' lean, in the head's own axes, and how
// fast it is changing: sideways, x, and forward, z, as parts of a turn.
type hairSwing struct{ u, v [2]float64 }

// swayed is the lean of the hair at each of n evenly spaced moments of the
// loop of length seconds that track follows, one cycle of the hair's own
// motion: a damped spring on the head's motion, and a light breeze. It
// runs the loop three times and keeps the last, so it comes back to where
// it began; a loop that is run through once runs it once, from rest.
func swayed(h *headTrack, length float64, n int, once bool) [][2]float64 {
	cycles := 3
	if once {
		cycles = 1
	}
	steps := int(math.Ceil(length / swayStep))
	dt := length / float64(steps)
	out := make([][2]float64, n)
	var s hairSwing
	prev, now := [3]float64{}, [3]float64{}
	for c := 0; c < cycles; c++ {
		for i := 0; i <= steps; i++ {
			t := dt * float64(i)
			p, r := h.at(t)
			if c == 0 && i == 0 {
				prev, now = p, p
			}
			a := scale(add(sub(p, scale(now, 2)), prev), 1/(dt*dt)) // the head's acceleration
			prev, now = now, p
			local := r.inverse().turn(a)
			s.step(local, t/length, dt)
			if c == cycles-1 {
				for k := range n {
					if j := int(math.Round(float64(k) / float64(max(n-1, 1)) * float64(steps))); j == i {
						out[k] = s.lean()
					}
				}
			}
		}
	}
	return out
}

// step moves the lean dt seconds on, the head accelerating by a in its
// own axes, the breeze at phase of the loop.
func (s *hairSwing) step(a [3]float64, phase, dt float64) {
	force := [2]float64{-a[0] / swayHang, -a[2] / swayHang}
	breeze := [2]float64{
		swayBreeze * (math.Sin(2*math.Pi*phase) + 0.5*math.Sin(2*math.Pi*2*phase+1.3)),
		swayBreeze * (math.Sin(2*math.Pi*phase+2.1) + 0.5*math.Sin(2*math.Pi*3*phase)),
	}
	for k := range 2 {
		acc := force[k] - swayRate*swayRate*(s.u[k]-breeze[k]) - 2*swayDamping*swayRate*s.v[k]
		s.v[k] += acc * dt
		s.u[k] += s.v[k] * dt
		if abs(s.u[k]) > swayLimit { // the hair meets the shoulders: it stops
			s.u[k], s.v[k] = math.Copysign(swayLimit, s.u[k]), 0
		}
	}
}

// lean is the hairSwing as two numbers.
func (s *hairSwing) lean() [2]float64 { return s.u }
