package farm

import (
	"math"
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// heading is a controller that keeps going one way, using what is held.
type heading vector.Vector3

func (h heading) control(*farm, *game.Object, float64) intent {
	return intent{dir: vector.Vector3(h), use: true}
}

// steer is a controller whose way can be changed as it goes.
type steer struct{ dir *vector.Vector3 }

func (s steer) control(*farm, *game.Object, float64) intent { return intent{dir: *s.dir} }

// The farmer wades into the pond, in the shallows by the shore, slower than
// they walk; dives in where the water is deep; swims there, slower still,
// shown sunk in the water; and, coming out, stands up and wades again.
func TestTheFarmerWadesDivesSwimsAndStandsUp(t *testing.T) {
	way := vector.Vector3{X: 1}
	g, f := newGame(steer{&way})
	o := g.Find("farmer")
	f.slot = 2 // the can, which would water a plot
	o.Transform.Position = standIn(f, 16, 8)
	var seen []string
	note := func() {
		a := o.Animator.Current()
		a = a[strings.LastIndex(a, "_")+1:]
		if a == "walk" && f.inWater(feet(o)) {
			a = "wade"
		}
		if len(seen) == 0 || seen[len(seen)-1] != a {
			seen = append(seen, a)
		}
	}
	// Into the pond and across it, until in its deep water.
	for range 600 {
		g.Step(frame)
		note()
		if f.deepAt(feet(o)) {
			break
		}
	}
	if !f.deepAt(feet(o)) {
		t.Fatalf("never reached the deep water; shown %v", seen)
	}
	for range 90 {
		g.Step(frame)
		note()
	}
	x := feet(o).X
	for range 30 {
		g.Step(frame)
	}
	if v := (feet(o).X - x) * 2; math.Abs(v-swimSpeed) > 1 {
		t.Errorf("swam at %v a second, not %v", v, swimSpeed)
	}
	// And back out.
	way.X = -1
	for range 600 {
		g.Step(frame)
		note()
		if !f.inWater(feet(o)) {
			break
		}
	}
	for range 120 {
		g.Step(frame)
		note()
	}
	want := []string{"walk", "wade", "dive", "swim", "rise", "wade", "walk"}
	if strings.Join(seen, " ") != strings.Join(want, " ") {
		t.Errorf("shown %v, not %v", seen, want)
	}
}

// Water next to land is shallow: a person swims only on water cells at
// least two cells from any land, the middle of the river or the lake,
// and wades nearer the shore.
func TestOneSwimsOnlyFarFromTheShore(t *testing.T) {
	_, f := newGame(idle{})
	deep, shallow := 0, 0
	for r := 0; r < f.ground.Rows; r++ {
		for c := 0; c < f.ground.Cols; c++ {
			p := standIn(f, c, r)
			if !f.deepAt(p) {
				if f.inWater(p) {
					shallow++
				}
				continue
			}
			deep++
			for dc := -1; dc <= 1; dc++ {
				for dr := -1; dr <= 1; dr++ {
					if !f.wet(c+dc, r+dr) {
						t.Fatalf("cell %d,%d swims beside the land at %d,%d", c, r, c+dc, r+dr)
					}
				}
			}
		}
	}
	if deep == 0 || shallow == 0 {
		t.Errorf("the pond has %d deep cells and %d shallow", deep, shallow)
	}
}

// The wild's lakes are water too, deep in the middle.
func TestTheWildsLakesAreWater(t *testing.T) {
	_, f := newGame(idle{})
	for r := -200; r < 200; r++ {
		for c := -200; c < 0; c++ {
			if f.wild.Gen.In("water", c, r) {
				if p := standIn(f, c, r); !f.inWater(p) {
					t.Fatalf("a lake's cell %d,%d is not water", c, r)
				}
				return
			}
		}
	}
	t.Fatal("no lake in the wild")
}
