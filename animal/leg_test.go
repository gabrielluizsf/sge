package animal

import (
	"math"
	"testing"
)

// A two-bar leg reaches the foot it is asked for, in the leg's plane,
// when the foot can be reached.
func TestALegReachesItsFoot(t *testing.T) {
	l := Leg{Hip: Vec{0, 2, 0}, Knee: Vec{0, 1, 0.5}, Foot: Vec{0, 0, 0}, Bend: 1, Upper: []int{0}, Lower: []int{1}}
	p := NewPose()
	to := Vec{0, 0, 0.6}
	l.Reach(p, to)
	// Walk the two bars from the hip with the turns Reach chose.
	up := turn(p.turns[0], sub(l.Knee, l.Hip))
	low := turn(p.turns[1], sub(l.Foot, l.Knee))
	foot := add(add(l.Hip, up), low)
	if dist(foot, to) > 1e-6 {
		t.Fatalf("the foot is at %v, not %v", foot, to)
	}
}

// Planted feet slide backward at the speed the gait says, and lifted ones
// come forward, so a body going at Speed does not slide them.
func TestAPlantedFootGoesBackAtTheBodysSpeed(t *testing.T) {
	g := Gait{Stride: 1.2, Period: 1, Duty: 0.6, Lift: 0.2}
	a, _ := g.offset(0.1)
	b, _ := g.offset(0.2)
	got := (a - b) / (0.1 * g.Period)
	if math.Abs(got-g.Speed()) > 1e-9 {
		t.Fatalf("a planted foot goes back at %v, the gait's speed is %v", got, g.Speed())
	}
	if _, up := g.offset(0.3); up != 0 {
		t.Fatalf("a planted foot is lifted by %v", up)
	}
	if _, up := g.offset(0.8); up <= 0 {
		t.Fatal("a swinging foot is not lifted")
	}
}
