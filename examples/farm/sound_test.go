package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// ears is a speaker that writes down what it plays.
type ears struct{ heard []string }

func (e *ears) play(name string, _ *vector.Vector3) { e.heard = append(e.heard, name) }

func (e *ears) count(name string) int {
	n := 0
	for _, h := range e.heard {
		if h == name {
			n++
		}
	}
	return n
}

// stroll is a controller that keeps going one way, using nothing.
type stroll vector.Vector3

func (s stroll) control(*farm, *game.Object) intent { return intent{dir: vector.Vector3(s)} }

// walk has the farmer walk the way dir points from cell col, row for n
// frames, and is what was heard.
func walk(dir vector.Vector3, from [2]int, n int) *ears {
	g, f := newGame(stroll(dir))
	e := &ears{}
	f.speaker = e
	g.Find("farmer").Transform.Position = standIn(f, from[0], from[1])
	for range n {
		g.Step(frame)
	}
	return e
}

// The farmer's steps sound of what they walk on, two to a stride, and
// swimming of the water: a splash going in, then strokes.
func TestStepsSoundOfTheGround(t *testing.T) {
	// Along the path, row 4; down column 1's soil; across the grass of
	// row 12; into the pond on row 8.
	for _, c := range []struct {
		dir  vector.Vector3
		from [2]int
		want string
	}{{vector.Vector3{X: 1}, [2]int{2, 4}, "step_path"}, {vector.Vector3{Y: 1}, [2]int{2, 6}, "step_soil"}, {vector.Vector3{X: 1}, [2]int{5, 12}, "step_grass"}} {
		e := walk(c.dir, c.from, 60)
		// Walking a second at 8 frames a second: four strides of two steps.
		if n := e.count(c.want); n < 3 || n > 5 || len(e.heard) != n {
			t.Errorf("walking from %v heard %v, not %s four times", c.from, e.heard, c.want)
		}
	}
	e := walk(vector.Vector3{X: 1}, [2]int{20, 8}, 120)
	if e.count("splash") != 1 || e.count("stroke") < 2 {
		t.Errorf("swimming heard %v", e.heard)
	}
}

// Buying, selling, sowing, watering and picking each have their sound,
// and what cannot be done its own.
func TestTheFarmsWorkIsHeard(t *testing.T) {
	_, f := newGame(idle{})
	e := &ears{}
	f.speaker = e
	f.money = 0
	f.choose(&shop{choice: 0})
	f.choose(&shop{choice: sellChoice})
	f.money = 20
	f.choose(&shop{choice: 0})
	f.harvest[0] = 1
	f.choose(&shop{choice: sellChoice})
	p := f.plots[0]
	f.slot = slotCarrot
	f.use(p.col, p.row)
	f.slot = slotCan
	f.use(p.col, p.row)
	p.stage = ripe
	f.use(p.col, p.row)
	if want := []string{"deny", "deny", "coin", "sell", "plant", "pour", "harvest"}; !slices.Equal(e.heard, want) {
		t.Errorf("heard %v, not %v", e.heard, want)
	}
}

// Near trees the ambience hears them, and now and then a bird in one;
// near the pond it hears water; in the open, nothing.
func TestTheAmbienceHearsTreesAndWater(t *testing.T) {
	g, f := newGame(idle{})
	if n := f.around(standIn(f, 25, 8)); n.water < 8 {
		t.Errorf("by the pond, %d cells of water", n.water)
	}
	if n := f.around(standIn(f, 19, 7)); len(n.trees) == 0 {
		t.Error("by a tree, no trees")
	}
	e := &ears{}
	f.speaker = e
	g.Sound = nil
	addSoundless(g, f)
	g.Find("farmer").Transform.Position = standIn(f, 19, 8)
	for range 60 * 20 {
		g.Step(frame)
	}
	if n := e.count("chirp"); n < 2 || n > 8 {
		t.Errorf("%d birds in 20 seconds by a tree", n)
	}
	if strings.Contains(strings.Join(e.heard, " "), "step") {
		t.Errorf("standing still heard %v", e.heard)
	}
}

// addSoundless adds the ambience to g with a mixer but keeps f's speaker.
func addSoundless(g *game.Game, f *farm) {
	s := f.speaker
	addSound(g, f)
	f.speaker = s
}
