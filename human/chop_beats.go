package human

import "math"

// beat is a pose a swing passes through, at a moment of the loop, from 0
// to 1, and how the body gets there from the beat before: the ease of its
// time, and whether it goes through the body in order.
type beat struct {
	at      float64
	pose    Pose
	ease    func(u float64) float64
	ordered bool
	// withTrunk has the left arm go with the trunk, not the right arm's
	// pace, as one that holds something still on the ground does.
	withTrunk bool
}

// hasten is an ease that gathers speed, as
// a swing falls; settle slows as it comes to rest.
func straight(u float64) float64 { return u }
func hasten(u float64) float64   { return math.Pow(u, 1.7) }
func settle(u float64) float64   { return 1 - (1-u)*(1-u) }

// poseAt is the pose of the beats at moment t of the loop: the segment t is
// in, eased from its beat to the next, which for the last is the first.
func poseAt(beats []beat, t float64) Pose {
	n := len(beats)
	for i := range n {
		next := beats[(i+1)%n]
		end := next.at
		if i == n-1 {
			end = 1
		}
		if t < beats[i].at || t >= end {
			continue
		}
		u := (t - beats[i].at) / (end - beats[i].at)
		if !next.ordered {
			return beats[i].pose.towards(next.pose, even(next.ease(u)))
		}
		r := reach{next.ease(min(u*1.35, 1)), next.ease(min(u*1.15, 1)), next.ease(u)}
		p := beats[i].pose.towards(next.pose, r)
		if next.withTrunk {
			p = withLeftArm(p, beats[i].pose.towards(next.pose, even(r.trunk)))
		}
		return p
	}
	return beats[0].pose
}

// withLeftArm is p with the joints of its left arm those of q.
func withLeftArm(p, q Pose) Pose {
	p.Shoulder[0], p.Raise[0], p.Plane[0], p.Elbow[0], p.Wrist[0] = q.Shoulder[0], q.Raise[0], q.Plane[0], q.Elbow[0], q.Wrist[0]
	return p
}
