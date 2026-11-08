package human

import (
	"math"
	"testing"
)

// The hammer is thrown with the right arm and the left only holds the
// iron: the right shoulder travels far, the left hardly moves.
func TestForgeIsTheRightArmsWork(t *testing.T) {
	lo, hi := [2]float64{9, 9}, [2]float64{-9, -9}
	for _, p := range Forge.Poses {
		for s := range 2 {
			lo[s], hi[s] = min(lo[s], p.Elbow[s]), max(hi[s], p.Elbow[s])
		}
	}
	if right := hi[1] - lo[1]; right < 1.5 {
		t.Errorf("the right elbow bends through only %.2f", right)
	}
	if left := hi[0] - lo[0]; left > 0.3 {
		t.Errorf("the left elbow bends through %.2f: it holds the iron", left)
	}
}

// A stroke a second: the hammer is raised slowly, over about half of it,
// thrown down fast, in about a fifth, and bounced up and rested after,
// the hips and trunk before the arm.
func TestForgeStrokeIsOneASecond(t *testing.T) {
	if Forge.Length != 1 {
		t.Errorf("a stroke takes %v seconds", Forge.Length)
	}
	top, hit := smithBeats[2], smithBeats[4]
	if top.at < 0.45 || top.at > 0.55 {
		t.Errorf("the hammer is raised by %.2f of the stroke", top.at)
	}
	if d := hit.at - top.at; d < 0.15 || d > 0.25 {
		t.Errorf("the blow takes %.2f of the stroke", d)
	}
	mid := poseAt(smithBeats, top.at+(hit.at-top.at)*0.3)
	legs := (mid.Hip[0] - top.pose.Hip[0]) / (hit.pose.Hip[0] - top.pose.Hip[0])
	arm := (mid.Elbow[1] - top.pose.Elbow[1]) / (hit.pose.Elbow[1] - top.pose.Elbow[1])
	if legs <= arm {
		t.Errorf("the hips have gone %.2f of the way and the arm %.2f: the arm leads", legs, arm)
	}
}

// strokes are the hammer's face and the iron's hot end, in the person's
// space, in each pose of Forge.
type stroke struct{ face, iron [3]float64 }

func strokes(t *testing.T) []stroke {
	t.Helper()
	if keptStrokes == nil {
		keptStrokes = makeStrokes(t)
	}
	return keptStrokes
}

// keptStrokes are the strokes made once for every test of them.
var keptStrokes []stroke

// makeStrokes makes the real forge's strokes.
func makeStrokes(t *testing.T) []stroke {
	t.Helper()
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Forge); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	var out []stroke
	for k := range Forge.Poses {
		hat, hq := handle(m, a, HammerGrip, k)
		tat, tq := handle(m, a, TongsGrip, k)
		out = append(out, stroke{
			face: add(hat, hq.turn([3]float64{0.15, 0.031, -0.05})),
			iron: add(tat, tq.turn([3]float64{0.2, 0.0125, 0})),
		})
	}
	return out
}

// At the blow the hammer's face is on the iron, which lies on the anvil's
// face, 83 cm up, ahead of the smith; at the top it is at the height of
// the head.
func TestForgeHammerMeetsTheIron(t *testing.T) {
	all := strokes(t)
	hit := int(smithBeats[4].at * float64(len(all)))
	s := all[hit]
	if d := sub(s.face, s.iron); math.Sqrt(dot3(d, d)) > 0.04 {
		t.Errorf("at the blow the face is %v from the iron", d)
	}
	if s.face[1] < 0.8 || s.face[1] > 0.92 {
		t.Errorf("the blow lands at %.2f m", s.face[1])
	}
	top := all[int(smithBeats[2].at*float64(len(all)))]
	if top.face[1] < 1.5 {
		t.Errorf("the hammer is raised only to %.2f m", top.face[1])
	}
	for k, s := range all {
		if s.face[1] < 0.8 {
			t.Errorf("pose %d: the hammer's face is %.2f m up, in the anvil", k, s.face[1])
		}
	}
}

// The left hand holds the iron on the anvil all through: it never leaves
// the anvil's face but to be lifted a little and turned.
func TestForgeLeftHandKeepsTheIronOnTheAnvil(t *testing.T) {
	for k, s := range strokes(t) {
		if s.iron[1] < 0.79 || s.iron[1] > 0.92 || math.Abs(s.iron[0]) > 0.06 || math.Abs(s.iron[2]-(0.52+anvilStep)) > 0.06 {
			t.Errorf("pose %d: the iron is at %.2f %.2f %.2f", k, s.iron[0], s.iron[1], s.iron[2])
		}
	}
}
