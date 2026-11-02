package human

import (
	"math"
	"testing"
)

// swing is the palms and the axe's place in a pose of the chop, as the
// real kit's person swings it.
type swing struct {
	left, right [3]float64
	at          [3]float64
	turn        quat
}

// swings are the real chop's swings, made once for every test of them.
func swings(t *testing.T) []swing { return swingsOf(t, Chop, ChopGrip) }

// swingsOf are the swings of motion mo and the grip g of its two hands,
// made once for each.
func swingsOf(t *testing.T, mo Motion, g Grip) []swing {
	t.Helper()
	if kept, ok := keptSwings[mo.Name]; ok {
		return kept
	}
	made := makeSwings(t, mo, g)
	if keptSwings == nil {
		keptSwings = map[string][]swing{}
	}
	keptSwings[mo.Name] = made
	return made
}

// keptSwings are the swings made, by the motion's name.
var keptSwings map[string][]swing

// makeSwings makes the swings of mo as the real kit's person holds the
// tool in g.
func makeSwings(t *testing.T, mo Motion, g Grip) []swing {
	t.Helper()
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, mo); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	var out []swing
	for k := range mo.Poses {
		w := worldAt(m, a, k)
		var s swing
		for i, h := range g.Hands {
			n := boneOf(m, h)
			p := add(w[n].at, w[n].turn.turn(palm(m, n)))
			if i == 0 {
				s.left = p
			} else {
				s.right = p
			}
		}
		s.at, s.turn = handle(m, a, g, k)
		out = append(out, s)
	}
	return out
}

// The left hand is on the handle in every pose, at its end, and the right
// is up it, a hand's breadth or two away, never torn off.
func TestChopHandsStayOnTheHandle(t *testing.T) {
	for k, s := range swings(t) {
		loc := s.turn.inverse().turn(sub(s.left, s.at))
		if loc[0] < -0.5 || loc[0] > 0.36 || abs(loc[1]-handleY) > 0.03 || abs(loc[2]) > 0.03 {
			t.Errorf("pose %d: the left palm is at %v on the axe", k, loc)
		}
		right := s.turn.inverse().turn(sub(s.right, s.at))
		if abs(right[1]-handleY) > 0.12 || abs(right[2]) > 0.12 || right[0] <= loc[0] || right[0] > 0.1 {
			t.Errorf("pose %d: the right palm is at %v on the axe", k, right)
		}
		if d := sub(s.right, s.left); dot3(d, d) < 0.08*0.08 || dot3(d, d) > 0.35*0.35 {
			t.Errorf("pose %d: the hands are %.2f apart", k, math.Sqrt(dot3(d, d)))
		}
	}
}

// At the ready the head is forward and up, at the trunk; at the strike
// the blade's edge faces the tree, forward along the facing and level, and
// never backwards or up while the axe comes down.
func TestChopBladeFacesTheTree(t *testing.T) {
	all := swings(t)
	head := func(k int) [3]float64 { return all[k].turn.turn([3]float64{1, 0, 0}) }
	blade := func(k int) [3]float64 { return all[k].turn.turn([3]float64{0, 0, 1}) }
	if h := head(0); h[1] < 0.3 || h[2] < 0.4 {
		t.Errorf("at the ready the head points %v", h)
	}
	hit := int(chopBeats[4].at * float64(len(all)))
	if b := blade(hit); b[2] < 0.9 || abs(b[1]) > 0.3 {
		t.Errorf("at the strike the edge faces %v", b)
	}
	if h := head(hit); h[1] > -0.4 {
		t.Errorf("at the strike the head points %v, not down", h)
	}
	for k := hit - 1; k <= hit+4; k++ {
		if b := blade(k); b[2] < 0.7 || b[1] > 0.4 {
			t.Errorf("pose %d: the edge faces %v", k, b)
		}
	}
}
