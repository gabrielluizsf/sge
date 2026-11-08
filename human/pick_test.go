package human

import "testing"

// pickings are the pose of picking, the palms' heights and the root's
// height, in each of the poses, as the real kit's person does it.
type picking struct {
	left, right [3]float64
	hips        float64
}

var pickings []picking

func picks(t *testing.T) []picking {
	t.Helper()
	if pickings != nil {
		return pickings
	}
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Pick); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	hand := func(w []spot, name string) [3]float64 {
		n := boneOf(m, name)
		return add(w[n].at, w[n].turn.turn(palm(m, n)))
	}
	for k := range Pick.Poses {
		w := worldAt(m, a, k)
		pickings = append(pickings, picking{hand(w, "wrist.L"), hand(w, "wrist.R"), w[boneOf(m, "root")].at[1]})
	}
	return pickings
}

// At the pull the hips are dropped into a squat and both hands are at the
// ground a stride ahead, on the plant, together; the back's bend is the
// hips' hinge, within what a back safely does, and the person ends
// standing, as they began.
func TestPickSquatsToTheGroundAndRisesStanding(t *testing.T) {
	all := picks(t)
	pull := int(pickBeats[3].at * float64(len(all)))
	if all[pull].hips > all[0].hips-0.3 {
		t.Errorf("the hips are at %.2f m, from %.2f standing", all[pull].hips, all[0].hips)
	}
	for _, h := range [][3]float64{all[pull].left, all[pull].right} {
		if h[1] > 0.3 || h[2] < 0.5 {
			t.Errorf("a hand is at %.2f m up, %.2f m ahead at the pull", h[1], h[2])
		}
	}
	if d := sub(all[pull].left, all[pull].right); dot3(d, d) > 0.25*0.25 {
		t.Errorf("the hands are %v apart on the plant", d)
	}
	last := all[len(all)-1]
	if last.hips < all[0].hips-0.05 || last.left[1] < 0.55 {
		t.Errorf("the person ends with hips at %.2f m, hands at %.2f m", last.hips, last.left[1])
	}
	for k, p := range Pick.Poses {
		if p.Lean > 0.95 {
			t.Errorf("pose %d: the back bends %.2f rad, past a safe hinge", k, p.Lean)
		}
	}
}

// The pull lifts the crop: the hands rise from the ground to the chest.
func TestPickLiftsTheCropToTheChest(t *testing.T) {
	all := picks(t)
	pull := int(pickBeats[3].at * float64(len(all)))
	top := 0.0
	for _, p := range all[pull:] {
		top = max(top, p.left[1])
	}
	if top < 0.95 {
		t.Errorf("the crop is lifted to %.2f m", top)
	}
}
