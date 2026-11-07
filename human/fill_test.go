package human

import "testing"

// The well is worked standing: the legs hardly bend, the left arm pulls
// the rope up and down twice while the right holds the can.
func TestWellIsWorkedStandingWithTwoPulls(t *testing.T) {
	var pulls int
	for k, p := range Well.Poses {
		if p.Knee[0] > 0.5 {
			t.Errorf("pose %d: the knee is bent %.2f at the well", k, p.Knee[0])
		}
		next := Well.Poses[(k+1)%len(Well.Poses)]
		if p.Elbow[0] < 1.2 && next.Elbow[0] >= 1.2 {
			pulls++
		}
	}
	if pulls != 2 {
		t.Errorf("%d pulls of the rope", pulls)
	}
	if WellFall.To <= WellFall.From || WellFall.Can.During[0] != Well.Name {
		t.Errorf("the water %v is not poured into the can held in %v", WellFall, WellFall.Can.During)
	}
}

// At the river the person crouches deep and low, as the well's never
// does, and rises again to stand.
func TestDipCrouchesAtTheBankAndRises(t *testing.T) {
	low, first, last := 0.0, Dip.Poses[0], Dip.Poses[len(Dip.Poses)-1]
	for _, p := range Dip.Poses {
		low = min(low, p.Bob)
	}
	if low > -0.5 {
		t.Errorf("the lowest the body goes is %.2f m", low)
	}
	if first.Knee[0] > 0.3 || last.Knee[0] > 0.3 {
		t.Errorf("the dip starts with a knee bent %.2f and ends with %.2f", first.Knee[0], last.Knee[0])
	}
	for _, p := range Well.Poses {
		if p.Bob < -0.1 {
			t.Fatalf("the well crouches %.2f m", p.Bob)
		}
	}
	if FillGrip.During[0] != Well.Name || FillGrip.During[1] != Dip.Name {
		t.Errorf("the can is held in %v", FillGrip.During)
	}
}

// The motions end by coming round to where they began.
func TestEatingAndFillingLoopCleanly(t *testing.T) {
	for _, m := range append(Eating, Drawing...) {
		a, b := m.Poses[len(m.Poses)-1], m.Poses[0]
		for _, d := range []float64{a.Nod - b.Nod, a.Lean - b.Lean, a.Elbow[1] - b.Elbow[1], a.Shoulder[1] - b.Shoulder[1]} {
			if d > 0.35 || d < -0.35 {
				t.Errorf("%s: the last pose is %.2f from the first", m.Name, d)
			}
		}
	}
}
