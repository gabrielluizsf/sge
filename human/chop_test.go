package human

import (
	"math"
	"testing"
)

// A felling cut is not a symmetric lift: the arms differ, the trunk turns
// over the trailing shoulder and back through the strike, and the weight
// moves from the back foot to the front.
func TestChopIsAsymmetric(t *testing.T) {
	var arms, twist [2]float64
	arms[0], twist[0] = math.Inf(1), math.Inf(1)
	twist[1], arms[1] = math.Inf(-1), 0
	for _, p := range Chop.Poses {
		arms[1] = max(arms[1], math.Abs(p.Shoulder[0]-p.Shoulder[1]))
		twist[0], twist[1] = min(twist[0], p.Twist), max(twist[1], p.Twist)
	}
	if arms[1] < 0.5 {
		t.Errorf("the arms differ by at most %.2f", arms[1])
	}
	if twist[1]-twist[0] < 0.8 {
		t.Errorf("the trunk turns through %.2f only", twist[1]-twist[0])
	}
	top, hit := chopBeats[2].pose, chopBeats[4].pose
	if top.Knee[1] <= top.Knee[0] || hit.Knee[0] <= hit.Knee[1] {
		t.Error("the weight does not move from the back foot to the front")
	}
}

// The backswing is slow, about half again the swing's time, and the strike
// is fast, a fifth of it, and the body begins it with the hips and trunk
// before the arms.
func TestChopBackswingIsSlowStrikeFast(t *testing.T) {
	top, hit := chopBeats[2], chopBeats[4]
	if top.at < 0.5 {
		t.Errorf("the backswing is done by %.2f of the swing", top.at)
	}
	if d := hit.at - top.at; d > 0.22 || d < 0.15 {
		t.Errorf("the strike takes %.2f of the swing", d)
	}
	mid := poseAt(chopBeats, top.at+(hit.at-top.at)*0.3)
	before := poseAt(chopBeats, top.at)
	legs := (mid.Hip[0] - before.Hip[0]) / (hit.pose.Hip[0] - before.Hip[0])
	arms := (mid.Elbow[0] - before.Elbow[0]) / (hit.pose.Elbow[0] - before.Elbow[0])
	if legs <= arms {
		t.Errorf("the hips have gone %.2f of the way and the arms %.2f: the arms lead", legs, arms)
	}
}

// A straight wrist is no turn, so no old motion is changed by having one;
// a bent wrist is a turn, and one past what a wrist can do is caught.
func TestWristBends(t *testing.T) {
	m := armed(t)
	nodes := map[string]int{}
	for i, n := range m.Nodes {
		nodes[n.Name] = i
	}
	a, err := restArms(m, nodes)
	if err != nil {
		t.Fatal(err)
	}
	if q := (Pose{}).turns(a)["wrist.R"]; q != still {
		t.Errorf("a straight wrist is turned %v", q)
	}
	if q := (Pose{Wrist: both(0.3)}).turns(a)["wrist.R"]; q == still {
		t.Error("a bent wrist is not turned")
	}
	if len((Pose{Wrist: both(2)}).Beyond()) != 2 {
		t.Error("a wrist bent two radians is healthy")
	}
}
