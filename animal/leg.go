package animal

import "math"

// Leg is two rigid bars of bones, the upper from the hip to the knee and
// the lower from the knee to the foot, moved by where the foot must be.
type Leg struct {
	Upper, Lower []int // bones of each bar, the first bone of each at its joint
	Hip, Knee    Vec   // where they rest
	Foot         Vec   // where the foot's tip rests
	Bend         float64
	Phase        float64
}

// At is the leg as it steps phase of a period after the others.
func (l Leg) At(bend, phase float64) Leg {
	l.Bend, l.Phase = bend, phase
	return l
}

// NewLeg is a leg of the bars' bones, whose foot ends at tip. bend is +1
// when the knee points forward, as a front leg's does, and -1 backward;
// phase is how far into the stride the leg is at the start, 0 to 1.
func (r *Rig) NewLeg(upper, lower []string, tip string) Leg {
	l := Leg{Upper: r.names(upper), Lower: r.names(lower)}
	l.Hip, l.Knee, l.Foot = r.Head(l.Upper[0]), r.Head(l.Lower[0]), r.Head(r.Bone(tip))
	return l
}

func (r *Rig) names(list []string) []int {
	out := make([]int, len(list))
	for i, n := range list {
		out[i] = r.Bone(n)
	}
	return out
}

// Reach turns the leg's bars so its foot is at the scene point to, with
// its hip where the pose's shift has put it: two-bar inverse kinematics in
// the leg's plane, forward and up.
func (l Leg) Reach(p *Pose, to Vec) {
	hip := add(l.Hip, p.Shift)
	l1 := plane(sub(l.Knee, l.Hip))
	l2 := plane(sub(l.Foot, l.Knee))
	n1, n2 := math.Hypot(l1[0], l1[1]), math.Hypot(l2[0], l2[1])
	t := plane(sub(to, hip))
	d := math.Hypot(t[0], t[1])
	d = min(max(d, math.Abs(n1-n2)+1e-6), (n1+n2)*0.999)
	psiT := math.Atan2(t[1], t[0])
	cosA := (n1*n1 + d*d - n2*n2) / (2 * n1 * d)
	a := math.Acos(min(max(cosA, -1), 1))
	psi1 := psiT - l.Bend*a
	knee := [2]float64{n1 * math.Cos(psi1), n1 * math.Sin(psi1)}
	foot := [2]float64{d * math.Cos(psiT), d * math.Sin(psiT)}
	psi2 := math.Atan2(foot[1]-knee[1], foot[0]-knee[0])
	q1 := About(Vec{1, 0, 0}, wrap(psi1-math.Atan2(l1[1], l1[0])))
	q2 := About(Vec{1, 0, 0}, wrap(psi2-math.Atan2(l2[1], l2[0])))
	for _, b := range l.Upper {
		p.Turn(b, q1)
	}
	for _, b := range l.Lower {
		p.Turn(b, q2)
	}
}

// plane is v in the leg's plane: Y first, then Z, so an angle grows as a
// turn about +X does.
func plane(v Vec) [2]float64 { return [2]float64{v[1], v[2]} }
