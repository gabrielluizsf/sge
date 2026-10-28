package scene

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// counter counts its updates.
type counter struct{ n int }

func (c *counter) Update(*game.Object, float64) { c.n++ }

// Going into a house leaves the world outside as it is, its objects
// still, and coming out finds it so; the player is in both.
func TestAScenesObjectsWaitWhileItIsLeft(t *testing.T) {
	g := game.New()
	player := g.Add(game.NewObject("player"))
	tree, ticks := game.NewObject("tree"), &counter{}
	tree.AddScript(ticks)
	g.Add(tree)
	g.Camera, g.Covered = vector.Vector3{X: 50}, true
	s := &Stage{Game: g, Keep: func(o *game.Object) bool { return o == player }, Current: "outside"}
	g.Step(1)
	built := 0
	inside := func(g *game.Game) {
		built++
		g.Add(game.NewObject("bed"))
		g.Camera, g.Covered, g.Background = vector.Vector3{}, false, canvas.RGB(1, 2, 3)
	}
	s.Switch("house", inside)
	g.Step(1)
	g.Step(1)
	if g.Find("tree") != nil || g.Find("bed") == nil || g.Find("player") == nil || ticks.n != 1 {
		t.Fatalf("inside: tree %v, bed %v, player %v, ticks %d", g.Find("tree"), g.Find("bed"), g.Find("player"), ticks.n)
	}
	if !s.Played("outside") || s.Current != "house" {
		t.Error("the world outside was not kept")
	}
	s.Switch("outside", nil)
	g.Step(1)
	if g.Find("tree") != tree || g.Find("bed") != nil || ticks.n != 2 || g.Camera.X != 50 || !g.Covered {
		t.Fatalf("outside again: ticks %d, camera %v, covered %v", ticks.n, g.Camera, g.Covered)
	}
	s.Switch("house", inside)
	g.Step(1)
	if built != 1 || g.Find("bed") == nil || g.Background != canvas.RGB(1, 2, 3) {
		t.Errorf("the house was built %d times", built)
	}
	s.Switch("house", inside)
	if built != 1 {
		t.Error("switching to the scene playing built it again")
	}
}

// A scene forgotten is built again.
func TestAForgottenSceneIsBuiltAgain(t *testing.T) {
	g := game.New()
	s := &Stage{Game: g, Current: "outside"}
	built := 0
	room := func(g *game.Game) { built++ }
	s.Switch("room", room)
	s.Switch("outside", nil)
	s.Forget("room")
	s.Switch("room", room)
	if built != 2 || s.Played("room") {
		t.Errorf("built %d times", built)
	}
}
