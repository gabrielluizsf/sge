package human

import (
	"fmt"
	"math"
)

// Range is what a joint may do, in radians: the least and the most.
type Range struct{ Min, Max float64 }

// deg is the range from min to max degrees.
func deg(min, max float64) Range { return Range{min * math.Pi / 180, max * math.Pi / 180} }

// Healthy is the range of motion of each of a Pose's joints in a healthy
// adult, from the clinical norms of Norkin and White, Measurement of Joint
// Motion, and the AAOS's (Greene and Heckman, 1994), each measured from
// the anatomical position and in the direction Pose says: hip flexion to
// 120 and extension to 20 degrees; knee flexion to 135; ankle
// plantarflexion to 50 and dorsiflexion to 20; shoulder flexion to 180
// and extension to 60, abduction to 180; elbow flexion to 150; the trunk
// bending to 80 and back 25, turning 45 each way; the neck bowing 45 and
// tilting back 70. A Plane is the shoulder's horizontal motion, 135
// forward and 45 back from the arm held out to the side, and a leg's
// Abduct 45 out and 30 across. Pitch and Bob are the body's place on the
// water or the ground, not a joint, and have no range.
var Healthy = map[string]Range{
	"Lean":     deg(-25, 80),
	"Twist":    deg(-45, 45),
	"Nod":      deg(-70, 45),
	"Hip":      deg(-20, 120),
	"Knee":     deg(0, 135),
	"Ankle":    deg(-20, 50),
	"Abduct":   deg(-30, 45),
	"Shoulder": deg(-60, 180),
	"Raise":    deg(-10, 180),
	"Plane":    deg(-45, 135),
	"Elbow":    deg(0, 150),
}

// Beyond lists every joint of p, by name and side, that is past its
// healthy range.
func (p Pose) Beyond() []string {
	var out []string
	check := func(label, name string, v float64) {
		if r := Healthy[name]; v < r.Min-1e-9 || v > r.Max+1e-9 {
			out = append(out, fmt.Sprintf("%s %.0f degrees, not %.0f to %.0f", label, v*180/math.Pi, r.Min*180/math.Pi, r.Max*180/math.Pi))
		}
	}
	check("Lean", "Lean", p.Lean)
	check("Twist", "Twist", p.Twist)
	check("Nod", "Nod", p.Nod)
	for s, side := range [2]string{"left", "right"} {
		for _, j := range []struct {
			name string
			v    [2]float64
		}{{"Hip", p.Hip}, {"Knee", p.Knee}, {"Ankle", p.Ankle}, {"Abduct", p.Abduct},
			{"Shoulder", p.Shoulder}, {"Raise", p.Raise}, {"Plane", p.Plane}, {"Elbow", p.Elbow}} {
			check(j.name+" "+side, j.name, j.v[s])
		}
	}
	return out
}
