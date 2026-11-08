package human

import (
	"math"
	"testing"
)

// The legs alternate: the right leg is the left half a stride on, and in
// each pose the legs are not the same.
func TestRunLegsAlternate(t *testing.T) {
	n := len(Run.Poses)
	different := 0
	for k, p := range Run.Poses {
		q := Run.Poses[(k+n/2)%n]
		if math.Abs(p.Hip[0]-q.Hip[1]) > 1e-9 || math.Abs(p.Knee[0]-q.Knee[1]) > 1e-9 {
			t.Fatalf("pose %d: the right leg is not the left half a stride on", k)
		}
		if math.Abs(p.Hip[0]-p.Hip[1]) > 0.3 {
			different++
		}
	}
	if different < n/2 {
		t.Errorf("the legs differ in only %d of %d poses", different, n)
	}
}

// Each arm swings against its own side's leg: forward as that leg is back.
func TestRunArmsSwingAgainstTheLegs(t *testing.T) {
	for k, p := range Run.Poses {
		for s := range 2 {
			if d := (p.Hip[s] - 0.2) * (p.Shoulder[s] - p.Shoulder[1-s]); d > 1e-9 && math.Abs(p.Hip[s]-p.Hip[1-s]) > 0.8 {
				t.Errorf("pose %d: the %d arm is with its own leg", k, s)
			}
		}
	}
}

// The knee is driven to about a radian of hip flexion, the heel comes up
// toward the buttock, the elbows are bent about a right angle, and the
// body leans a little from the ankles.
func TestRunStridePosture(t *testing.T) {
	hip, knee := 0.0, 0.0
	var elbow float64
	for _, p := range Run.Poses {
		hip, knee = max(hip, p.Hip[0]), max(knee, p.Knee[0])
		elbow += p.Elbow[0] + p.Elbow[1]
		if p.Lean < 0.08 || p.Lean > 0.3 {
			t.Errorf("the lean is %.2f", p.Lean)
		}
	}
	if hip < 0.9 || hip > 1.15 {
		t.Errorf("the knee is driven to %.2f radians of hip flexion", hip)
	}
	if knee < 1.6 {
		t.Errorf("the heel comes up with the knee at %.2f radians", knee)
	}
	if mean := elbow / float64(2*len(Run.Poses)); mean < 1.5 || mean > 1.9 {
		t.Errorf("the elbows are bent %.2f radians on average", mean)
	}
}

// Running has a flight phase, twice a stride, when both feet are off the
// ground; and otherwise a foot is on it; the body bobs a little.
func TestRunHasAFlightPhase(t *testing.T) {
	flying, lo, hi := 0, 9.0, -9.0
	for k, p := range Run.Poses {
		left, right := footClearance(p, 0), footClearance(p, 1)
		if min(left, right) > 0.01 {
			flying++
		} else if min(left, right) > 0.03 {
			t.Errorf("pose %d: neither foot is on the ground, nor quite off it", k)
		}
		lo, hi = min(lo, p.Bob), max(hi, p.Bob)
	}
	if flying < 2 || flying > 6 {
		t.Errorf("both feet are off the ground in %d of %d poses", flying, len(Run.Poses))
	}
	if hi-lo < 0.03 || hi-lo > 0.12 {
		t.Errorf("the body bobs %.2f m", hi-lo)
	}
	if Run.Length < 0.65 || Run.Length > 0.75 {
		t.Errorf("a stride pair takes %v seconds", Run.Length)
	}
}
