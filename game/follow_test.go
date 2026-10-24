package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/vector"
)

// Follow keeps its object in the middle of the view, whatever its size.
func TestFollowCentresTheCameraOnItsObject(t *testing.T) {
	g := New()
	o := g.Add(NewObject("hero"))
	o.AddScript(Follow{})
	o.Transform.Position = vector.Vector3{X: 300, Y: 200}
	cv, _ := canvas.NewCanvas(100, 60)
	g.Draw(cv) // the view's size is known once drawn
	g.Step(0.1)
	if g.Camera != (vector.Vector3{X: 250, Y: 170}) {
		t.Errorf("camera at %v, not centred on 300,200", g.Camera)
	}
	big, _ := canvas.NewCanvas(200, 100)
	g.Draw(big)
	g.Step(0.1)
	if p := g.ToScreen(o.Transform.Position); p.X != 100 || p.Y != 50 {
		t.Errorf("after the window grew, the object is drawn at %v, not in the middle", p)
	}
}

// A lagging camera eases towards its object rather than jumping.
func TestFollowWithLagEasesTowardsItsObject(t *testing.T) {
	g := New()
	o := g.Add(NewObject("hero"))
	o.AddScript(Follow{Lag: 0.2})
	o.Transform.Position = vector.Vector3{X: -300}
	cv, _ := canvas.NewCanvas(100, 100)
	g.Draw(cv)
	g.Step(0.016)
	if g.Camera.X != -350 {
		t.Fatalf("the camera starts at %v, not on its object", g.Camera.X)
	}
	o.Transform.Position = vector.Vector3{X: 100}
	g.Step(0.016)
	if g.Camera.X <= -350 || g.Camera.X >= -100 {
		t.Errorf("camera jumped to %v", g.Camera.X)
	}
	for range 200 {
		g.Step(0.016)
	}
	if d := g.Camera.X - 50; d > 0.5 || d < -0.5 {
		t.Errorf("camera settled at %v, not 50", g.Camera.X)
	}
}
