package human

import "testing"

// The hoe's blade, in the hoe's own frame, is the face at the handle's
// far end: its middle, where it bites.
var hoeBlade = [3]float64{0.48, 0.0275, 0}

// blade is where the hoe's blade is in a swing, and the way along the
// handle toward it.
func (s swing) blade() (at, along [3]float64) {
	return add(s.at, s.turn.turn(hoeBlade)), s.turn.turn([3]float64{1, 0, 0})
}

// Both hands are on the handle in every pose, the left near its end and
// the right lower down it, a forearm's length or less from the left, and
// the blade never goes into the ground.
func TestHoeHandsStayOnTheHandle(t *testing.T) {
	for k, s := range swingsOf(t, Hoe, HoeGrip) {
		left := s.turn.inverse().turn(sub(s.left, s.at))
		right := s.turn.inverse().turn(sub(s.right, s.at))
		if left[0] > -0.46 || left[0] < -0.56 || abs(left[1]-handleY) > 0.03 || abs(left[2]) > 0.03 {
			t.Errorf("pose %d: the left palm is at %v on the hoe", k, left)
		}
		if right[0] <= left[0]+0.06 || right[0] > 0.1 || abs(right[1]-handleY) > 0.12 || abs(right[2]) > 0.12 {
			t.Errorf("pose %d: the right palm is at %v on the hoe", k, right)
		}
		if at, _ := s.blade(); at[1] < -0.03 {
			t.Errorf("pose %d: the blade is %.2f m into the ground", k, -at[1])
		}
	}
}

// At the bite the blade is at the ground, a stride ahead of the feet, its
// edge facing down and forward; before it, it was lifted to the height of
// the shoulders; and after it the stroke draws the blade toward the body.
func TestHoeBitesAStrideAheadAndDrawsBack(t *testing.T) {
	all := swingsOf(t, Hoe, HoeGrip)
	bite := int(hoeBeats[4].at * float64(len(all)))
	at, along := all[bite].blade()
	if at[1] > 0.06 || at[2] < 0.75 {
		t.Errorf("the blade bites at %.2f m up, %.2f m ahead", at[1], at[2])
	}
	if along[1] > -0.6 || along[2] < 0.2 {
		t.Errorf("the edge faces %v at the bite, not down and forward", along)
	}
	var top float64
	for _, s := range all[:bite] {
		b, _ := s.blade()
		top = max(top, b[1])
	}
	if top < 1.2 {
		t.Errorf("the blade is lifted to only %.2f m", top)
	}
	end, _ := all[int(hoeBeats[6].at*float64(len(all)))].blade()
	if drawn := at[2] - end[2]; drawn < 0.2 || end[1] > 0.1 {
		t.Errorf("the blade is drawn %.2f m back, to %.2f m up", drawn, end[1])
	}
	prev := at[2]
	for k := bite + 1; k <= int(hoeBeats[6].at*float64(len(all))); k++ {
		b, _ := all[k].blade()
		if b[2] > prev+0.02 {
			t.Errorf("pose %d: the blade goes forward in the draw, to %.2f m from %.2f", k, b[2], prev)
		}
		prev = b[2]
	}
}
