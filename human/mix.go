package human

// lerp is a to b, f of the way.
func lerp(a, b, f float64) float64 { return a + (b-a)*f }

// lerp2 is each side's lerp.
func lerp2(a, b [2]float64, f float64) [2]float64 {
	return [2]float64{lerp(a[0], b[0], f), lerp(a[1], b[1], f)}
}

// mix is the pose f of the way from a to b, joint by joint.
func mix(a, b Pose, f float64) Pose {
	return Pose{
		Pitch: lerp(a.Pitch, b.Pitch, f), Lean: lerp(a.Lean, b.Lean, f), Twist: lerp(a.Twist, b.Twist, f), Nod: lerp(a.Nod, b.Nod, f),
		Hip: lerp2(a.Hip, b.Hip, f), Knee: lerp2(a.Knee, b.Knee, f), Ankle: lerp2(a.Ankle, b.Ankle, f), Abduct: lerp2(a.Abduct, b.Abduct, f),
		Shoulder: lerp2(a.Shoulder, b.Shoulder, f), Raise: lerp2(a.Raise, b.Raise, f), Elbow: lerp2(a.Elbow, b.Elbow, f),
		Plane: lerp2(a.Plane, b.Plane, f), Bob: lerp(a.Bob, b.Bob, f),
	}
}

// ease is f eased in and out, from 0 to 1.
func ease(f float64) float64 { return f * f * (3 - 2*f) }
