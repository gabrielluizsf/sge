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

// The farmer walks into the pond and swims in it, slower than they walk,
// shown sunk in the water.
func TestTheFarmerSwimsInThePond(t *testing.T) {
	g, f := newGame(heading{X: 1})
	o := g.Find("farmer")
	f.slot = 2 // the can, which would water a plot
	o.Transform.Position = standIn(f, 16, 8)
	for range 150 {
		g.Step(frame)
	}
	if col, row, _ := f.ground.CellAt(feet(o).X, feet(o).Y); !f.inWater(feet(o)) || col < 22 {
		t.Fatalf("at cell %d,%d, not in the pond", col, row)
	}
	if a := o.Animator.Current(); !strings.HasSuffix(a, "_swim") {
		t.Errorf("shown %q in the water", a)
	}
	x := feet(o).X
	for range 30 {
		g.Step(frame)
	}
	if v := (feet(o).X - x) * 2; math.Abs(v-swimSpeed) > 1 {
		t.Errorf("swam at %v a second, not %v", v, swimSpeed)
	}
	if !f.inWater(feet(o)) {
		t.Fatal("swam out of the pond")
	}
}

// The wild's lakes are water too.
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
