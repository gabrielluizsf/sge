package human

import "testing"

// pouring is a pose of the watering, as the real kit's person does it:
// where the can is, how it is turned, and where the hand it is held in
// has its palm.
type pouring struct {
	at, palm [3]float64
	turn     quat
}

// tip is where the can's spout ends, and out which way it points.
func (p pouring) tip() [3]float64 { return add(p.at, p.turn.turn(canSpout)) }
func (p pouring) out() [3]float64 { return p.turn.turn(normalize(canOut)) }

// pourings are the poses of the real watering, made once.
var pourings []pouring

func pours(t *testing.T) []pouring {
	t.Helper()
	if pourings != nil {
		return pourings
	}
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Water); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	n := boneOf(m, CanGrip.Hands[0])
	for k := range Water.Poses {
		w := worldAt(m, a, k)[n]
		at, q := handle(m, a, CanGrip, k)
		pourings = append(pourings, pouring{at, add(w.at, w.turn.turn(palm(m, n))), q})
	}
	return pourings
}

// The hand holds the can's handle in every pose: the palm is where the
// handle is.
func TestWaterHandHoldsTheHandleAlways(t *testing.T) {
	for k, p := range pours(t) {
		grip := add(p.at, p.turn.turn(canHandle))
		if d := sub(p.palm, grip); dot3(d, d) > 0.02*0.02 {
			t.Errorf("pose %d: the palm is %v from the handle", k, d)
		}
	}
}

// While it pours, the can is tipped so that the spout points down at the
// soil a stride ahead, low, and lower than at any other time; the pour is
// a third to a half of the watering.
func TestWaterSpoutPointsAtTheSoilWhilePouring(t *testing.T) {
	all := pours(t)
	var pouringKeys int
	lowest, lowestAt, other := 9.0, 0, 9.0
	for k, p := range all {
		if p.out()[1] < -0.3 {
			pouringKeys++
			if tip := p.tip(); tip[1] < lowest {
				lowest, lowestAt = tip[1], k
			}
		}
	}
	for _, p := range all {
		if p.out()[1] >= -0.3 {
			other = min(other, p.tip()[1])
		}
	}
	if share := float64(pouringKeys) / float64(len(all)); share < 0.3 || share > 0.55 {
		t.Errorf("the can pours through %.2f of the watering", share)
	}
	tip := all[lowestAt].tip()
	if tip[1] < 0.3 || tip[1] > 0.75 || tip[2] < 0.6 {
		t.Errorf("the spout pours from %.2f m up, %.2f m ahead", tip[1], tip[2])
	}
	if lowest > other {
		t.Errorf("the spout is at %.2f m when pouring, %.2f at its lowest otherwise", lowest, other)
	}
}

// Not pouring, the can is upright: carried at the start and the end.
func TestWaterCanIsUprightWhenNotPouring(t *testing.T) {
	all := pours(t)
	for _, k := range []int{0, len(all) - 1} {
		if up := all[k].turn.turn([3]float64{0, 1, 0}); up[1] < 0.9 {
			t.Errorf("pose %d: the can leans, its top pointing %v", k, up)
		}
		if all[k].out()[1] < -0.3 {
			t.Errorf("pose %d: the can pours", k)
		}
	}
}

// The watering takes two seconds and a pour is a whole of it: it ends as
// it began, so it loops.
func TestWaterTakesTwoSeconds(t *testing.T) {
	if Water.Length != 2 {
		t.Errorf("the watering takes %v seconds", Water.Length)
	}
}
