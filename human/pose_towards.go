package human

// reach is how far a pose has gone toward another, from 0 to 1, as the
// body's parts go in turn: a swing starts in the legs and hips, passes
// through the trunk and ends in the arms.
type reach struct{ legs, trunk, arms float64 }

// even is the same reach for every part.
func even(f float64) reach { return reach{f, f, f} }

// towards is p, r of the way to q: its legs, its trunk and head, and its
// arms each as far as r says.
func (p Pose) towards(q Pose, r reach) Pose {
	return Pose{
		Pitch: lerp(p.Pitch, q.Pitch, r.trunk), Lean: lerp(p.Lean, q.Lean, r.trunk),
		Twist: lerp(p.Twist, q.Twist, r.trunk), Nod: lerp(p.Nod, q.Nod, r.trunk),
		Hip: lerp2(p.Hip, q.Hip, r.legs), Knee: lerp2(p.Knee, q.Knee, r.legs),
		Ankle: lerp2(p.Ankle, q.Ankle, r.legs), Abduct: lerp2(p.Abduct, q.Abduct, r.legs),
		Shoulder: lerp2(p.Shoulder, q.Shoulder, r.arms), Raise: lerp2(p.Raise, q.Raise, r.arms),
		Elbow: lerp2(p.Elbow, q.Elbow, r.arms), Plane: lerp2(p.Plane, q.Plane, r.arms),
		Wrist: lerp2(p.Wrist, q.Wrist, r.arms),
		Bob:   lerp(p.Bob, q.Bob, r.legs), Forward: lerp(p.Forward, q.Forward, r.legs),
	}
}
