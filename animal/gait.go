package animal

import "math"

// Gait is how a walk goes: each foot steps Stride metres along, down on
// the ground for Duty of each Period seconds and lifted Lift metres for
// the rest.
type Gait struct {
	Stride, Period, Duty, Lift float64
}

// Speed is how fast the body goes, in metres a second, for its planted
// feet not to slide.
func (g Gait) Speed() float64 { return g.Stride / (g.Duty * g.Period) }

// offset is where a foot is from its rest place, forward and up, at
// u, the fraction of the period it is into its stride.
func (g Gait) offset(u float64) (dz, dy float64) {
	u -= math.Floor(u)
	if u < g.Duty {
		return g.Stride/2 - g.Stride*u/g.Duty, 0
	}
	s := (u - g.Duty) / (1 - g.Duty)
	ease := s * s * (3 - 2*s)
	return -g.Stride/2 + g.Stride*ease, g.Lift * math.Sin(math.Pi*s)
}

// Walk is the pose of the legs at second t of a walk, the body lowered by
// crouch: planted feet go backward under the body, lifted ones forward.
func Walk(legs []Leg, g Gait, t, crouch float64) *Pose {
	p := NewPose()
	p.Shift = Vec{0, -crouch, 0}
	for _, l := range legs {
		dz, dy := g.offset(t/g.Period + l.Phase)
		l.Reach(p, add(l.Foot, Vec{0, dy, dz}))
	}
	return p
}
