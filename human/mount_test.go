package human

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// climbBody is the real kit's person with the motions of getting on and down.
func climbBody(t *testing.T) *gltf.Model {
	t.Helper()
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Mount, Dismount); err != nil {
		t.Fatal(err)
	}
	return m
}

// distance is how far apart two spots are.
func distance(a, b [3]float64) float64 {
	d := [3]float64{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
	return math.Sqrt(d[0]*d[0] + d[1]*d[1] + d[2]*d[2])
}

// timeAt is how far through a climb frame k is.
func timeAt(k int) float64 { return float64(k) / (mountFrames - 1) }

// A climb starts and ends in the very poses on either side of it.
func TestAClimbMeetsTheMotionsAroundIt(t *testing.T) {
	if Mount.Poses[0] != standBeside() || Mount.Poses[len(Mount.Poses)-1] != ridePose(0) {
		t.Error("getting on does not go from standing to the seat")
	}
	if Dismount.Poses[0] != ridePose(0) || Dismount.Poses[len(Dismount.Poses)-1] != standOut() {
		t.Error("getting down does not go from the seat to standing")
	}
}

// The foot in the stirrup stays there, from when it is put in until the seat or
// the slide down, and the one on the floor stays there until the push.
func TestTheFeetStayWhereTheyBear(t *testing.T) {
	m := climbBody(t)
	cases := []struct {
		name      string
		anim      int
		bone      string
		spot      [3]float64
		from, til float64
	}{
		{"getting on, in the stirrup", 0, "foot.R", nearStirrup, 0.28, 1},
		{"getting on, on the floor", 0, "foot.L", backOnFloor, 0.16, 0.34},
	}
	for _, c := range cases {
		for k := range mountFrames {
			if u := timeAt(k); u < c.from || u > c.til {
				continue
			}
			w := worldAt(m, &m.Animations[c.anim], k)
			if d := distance(w[boneOf(m, c.bone)].at, c.spot); d > 0.09 {
				t.Errorf("%s: frame %d is %.3f from its spot", c.name, k, d)
			}
		}
	}
}

// The leg that goes over goes clear of the horse's back, though a knee may press
// its side.
func TestTheSwingingLegClearsTheHorse(t *testing.T) {
	m := climbBody(t)
	for k := 1; k < mountFrames-1; k++ {
		w := worldAt(m, &m.Animations[0], k)
		for _, bone := range []string{"foot.L", "lowerleg01.L"} {
			at := w[boneOf(m, bone)].at
			round := math.Hypot(at[0], at[1]-BarrelHeight)
			if math.Abs(at[2]) < BarrelLength/2 && round < BarrelRadius-0.04 {
				t.Errorf("frame %d: %s inside the barrel, %.3f from its axis", k, bone, round)
			}
		}
	}
}

// The hands are on the horse from when they reach it until the seat.
func TestTheHandsHoldTheHorse(t *testing.T) {
	m := climbBody(t)
	pommel, mane := [3]float64{-0.20, 1.44, 0.06}, [3]float64{-0.17, 1.50, 0.38}
	for k := range mountFrames {
		if u := timeAt(k); u < 0.16 || u > 0.85 {
			continue
		}
		w := worldAt(m, &m.Animations[0], k)
		for _, h := range []struct {
			bone string
			spot [3]float64
		}{{"wrist.L", pommel}, {"wrist.R", mane}} {
			n := boneOf(m, h.bone)
			at := add(w[n].at, w[n].turn.turn(palm(m, n)))
			if d := distance(at, h.spot); d > 0.2 {
				t.Errorf("frame %d: %s %.3f from where it holds", k, h.bone, d)
			}
		}
	}
}

// No joint goes past what a body does, and nothing jumps from frame to frame.
func TestAClimbIsABodysMotion(t *testing.T) {
	var last [3]float64
	for k, p := range Mount.Poses {
		for s := range 2 {
			if p.Hip[s] > 2.45 || p.Knee[s] > 2.4 || math.Abs(p.Abduct[s]) > math.Pi/2+1e-9 {
				t.Errorf("frame %d leg %d: hip %.2f knee %.2f out %.2f", k, s, p.Hip[s], p.Knee[s], p.Abduct[s])
			}
		}
		at := rootAt(p)
		if k > 0 && distance(at, last) > 0.28 {
			t.Errorf("frame %d: the body jumps %.3f", k, distance(at, last))
		}
		last = at
	}
}

// A curve goes through its keys, stays still along equal ones and goes beyond
// none of them.
func TestACurveDoesNotOvershoot(t *testing.T) {
	c := curve{at: []float64{0, 0.3, 0.6, 0.7, 1}, keys: []float64{0, 5, 5, 1, 1}}
	for i, at := range c.at {
		if got := c.of(at); math.Abs(got-c.keys[i]) > 1e-9 {
			t.Errorf("at %v: %v, not its key %v", at, got, c.keys[i])
		}
	}
	for u := 0.0; u <= 1; u += 0.01 {
		v := c.of(u)
		if v < -1e-9 || v > 5+1e-9 {
			t.Errorf("at %v: %v, beyond its keys", u, v)
		}
		if u >= 0.3 && u <= 0.6 && math.Abs(v-5) > 1e-9 {
			t.Errorf("at %v: %v moves between equal keys", u, v)
		}
	}
}

// A leg sent to a spot of its reach puts the real body's ankle there, turned or
// not.
func TestALegReachesItsSpot(t *testing.T) {
	spots := map[float64][][3]float64{
		0:   {{-0.76, 0.5, 0.3}, {-0.8, 0.091, -0.1}},
		1.3: {{-0.279, 0.814, -0.065}, {-0.6, 0.091, 0}, {-0.5, 0.6, 0}},
	}
	for turn, list := range spots {
		p := Stand.Poses[0]
		p.Turn, p.Aside, p.Forward = turn, -0.65, -0.1
		for _, spot := range list {
			q := p
			legTo(&q, stirrupLeg, spot)
			mo := Motion{Name: "reach", Length: 1, Poses: []Pose{q, q}, Once: true}
			m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
			if err != nil {
				t.Fatal(err)
			}
			if err := Animate(m, mo); err != nil {
				t.Fatal(err)
			}
			w := worldAt(m, &m.Animations[0], 0)
			if d := distance(w[boneOf(m, "foot.R")].at, spot); d > 0.09 {
				t.Errorf("turned %v, to %v: %.3f away", turn, spot, d)
			}
		}
	}
}

// Getting down sits sideways first: both feet hanging on the side one gets down
// by, the far one having come over; then both are off the saddle and clear of
// the floor, in a jump; then down on bent knees.
func TestGettingDownSitsSidewaysAndJumps(t *testing.T) {
	m := climbBody(t)
	at := func(k int, bone string) [3]float64 {
		return worldAt(m, &m.Animations[1], k)[boneOf(m, bone)].at
	}
	seat, jump, land := -1, -1, -1
	for k := range mountFrames {
		u := timeAt(k)
		if seat < 0 && u >= 0.52 {
			seat = k
		}
		if u <= 0.74 {
			jump = k
		}
		if land < 0 && u >= 0.86 {
			land = k
		}
	}
	for _, bone := range []string{"foot.L", "foot.R"} {
		if p := at(seat, bone); p[0] > -0.35 {
			t.Errorf("sitting sideways: %s at %.2f, not on the near side", bone, p[0])
		}
		if p := at(jump, bone); p[1] < 0.3 {
			t.Errorf("in the jump: %s only %.2f up", bone, p[1])
		}
		if p := at(land, bone); p[1] > 0.14 {
			t.Errorf("landing: %s %.2f up", bone, p[1])
		}
	}
	if p := Dismount.Poses[land]; p.Knee[0] < 0.9 || p.Knee[1] < 0.9 {
		t.Errorf("landing with knees at %.2f, %.2f", p.Knee[0], p.Knee[1])
	}
	if !(rootAt(Dismount.Poses[land])[0] < -0.5) {
		t.Errorf("landing %.2f from the horse's middle", rootAt(Dismount.Poses[land])[0])
	}
}

// The body at the end of getting down is the same as standing, a quarter turned,
// on the spot beside the horse: the picture does not jump when the game swaps it.
func TestGettingDownEndsWhereStandingIs(t *testing.T) {
	p := standOut()
	r := rootAt(p)
	if math.Abs(r[0]-(-MountAside+0.065)) > 0.002 || math.Abs(r[2]) > 0.002 {
		t.Errorf("the root at %v", r)
	}
}

// No joint of getting down goes past what a body does.
func TestGettingDownIsABodysMotion(t *testing.T) {
	for k, p := range Dismount.Poses {
		for s := range 2 {
			if p.Hip[s] > 2.45 || p.Knee[s] > 2.4 || math.Abs(p.Abduct[s]) > math.Pi/2+1e-9 {
				t.Errorf("frame %d leg %d: hip %.2f knee %.2f out %.2f", k, s, p.Hip[s], p.Knee[s], p.Abduct[s])
			}
		}
	}
}
