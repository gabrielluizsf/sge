package game_test

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// square is an object showing a square of colour c at p.
func square(c canvas.Color, p vector.Vector3) *game.Object {
	img, _ := canvas.NewCanvas(8, 8)
	img.Clear(c)
	o := game.NewObject("square")
	o.Sprite = &sprite.Sprite{Image: img, PivotX: 0.5, PivotY: 0.5}
	o.Transform.Position = p
	return o
}

// In a projected world, of two objects of one Z the nearer is drawn over
// the farther, whichever was added first.
func TestProjectedObjectsDrawNearerOverFarther(t *testing.T) {
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	red, blue := canvas.RGB(255, 0, 0), canvas.RGB(0, 0, 255)
	g.Add(square(red, vector.Vector3{X: 40, Y: 40}))  // nearer
	g.Add(square(blue, vector.Vector3{X: 38, Y: 38})) // farther, overlapping it
	cv, _ := canvas.NewCanvas(200, 200)
	g.Draw(cv)
	at := g.ToScreen(vector.Vector3{X: 40, Y: 40})
	if got := cv.At(int(at.X), int(at.Y)); got != red {
		t.Errorf("the nearer square is under the farther: %v", got)
	}
}

// Follow centres its object on the screen in a projected world too, and
// ToWorld finds the ground under a point of the screen.
func TestProjectedFollowAndToWorld(t *testing.T) {
	g := game.New()
	g.Projection = iso.New(32, 64, 32)
	o := g.Add(game.NewObject("hero"))
	o.AddScript(game.Follow{})
	o.Transform.Position = vector.Vector3{X: 300, Y: 200}
	cv, _ := canvas.NewCanvas(100, 60)
	g.Draw(cv)
	g.Step(0.1)
	if p := g.ToScreen(o.Transform.Position); p.Sub(vector.Vector3{X: 50, Y: 30}).Length() > 1 {
		t.Errorf("the object is drawn at %v, not in the middle", p)
	}
	if w := g.ToWorld(vector.Vector3{X: 50, Y: 30}); w.Sub(o.Transform.Position).Length() > 2 {
		t.Errorf("the ground under the middle is %v, not the object's %v", w, o.Transform.Position)
	}
	v := g.VisibleWorld()
	if p := o.Transform.Position; p.X < v.X || p.Y < v.Y || p.X > v.X+v.Width || p.Y > v.Y+v.Height {
		t.Errorf("the visible world %v leaves out the object at %v", v, p)
	}
}
