package human

import "testing"

// The right hand reaches the mouth at the bite and is back at the chest
// for the chewing, the head nods a few times while it chews, and the left
// arm stays loose all along.
func TestEatBringsTheFoodToTheMouthAndChews(t *testing.T) {
	var up, low, chews int
	prev := 0.0
	for k, p := range Eat.Poses {
		if p.Elbow[1] > 2.0 {
			up++
		}
		if p.Elbow[1] < 1.7 && k > len(Eat.Poses)/2 {
			low++
		}
		if k > 0 && k > len(Eat.Poses)/2 && (p.Nod-prev)*(prev-Eat.Poses[k-2%len(Eat.Poses)].Nod) < 0 {
			chews++
		}
		prev = p.Nod
		if p.Elbow[0] > 0.6 || p.Shoulder[0] > 0.3 {
			t.Errorf("pose %d: the left arm works: elbow %.2f shoulder %.2f", k, p.Elbow[0], p.Shoulder[0])
		}
	}
	if up == 0 || low == 0 || chews < 3 {
		t.Errorf("hand up in %d poses, down in %d, %d chews", up, low, chews)
	}
	if len(MorselGrip.During) != 1 || MorselGrip.During[0] != Eat.Name {
		t.Errorf("the food is held in %v", MorselGrip.During)
	}
}

// At the bite the hand is at the mouth, as the real person has it.
func TestEatHandReachesTheMouth(t *testing.T) {
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Eat); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	wr, jaw := boneOf(m, "wrist.R"), boneOf(m, "jaw")
	best := 9.0
	for k := range Eat.Poses {
		w := worldAt(m, a, k)
		hand := add(w[wr].at, w[wr].turn.turn(palm(m, wr)))
		d := sub(hand, w[jaw].at)
		best = min(best, dot3(d, d))
	}
	if best > 0.14*0.14 {
		t.Errorf("the hand comes no nearer the jaw than %.2f m", best)
	}
}
