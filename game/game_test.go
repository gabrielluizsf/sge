package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// recorder is a script that writes down what happens to its object.
type recorder struct {
	log *[]string
}

func (r recorder) Start(o *Object)              { *r.log = append(*r.log, "start "+o.Name) }
func (r recorder) Update(o *Object, dt float64) { *r.log = append(*r.log, "update "+o.Name) }
func (r recorder) Collision(o, other *Object)   { *r.log = append(*r.log, o.Name+" hit "+other.Name) }

func TestScriptsStartOnceThenUpdate(t *testing.T) {
	var log []string
	g := New()
	g.Add(NewObject("a")).AddScript(recorder{&log})
	g.Step(0.1)
	g.Step(0.1)
	want := []string{"start a", "update a", "update a"}
	if len(log) != len(want) {
		t.Fatalf("log = %v, want %v", log, want)
	}
	for i := range want {
		if log[i] != want[i] {
			t.Fatalf("log = %v, want %v", log, want)
		}
	}
	if g.Time() < 0.2-1e-9 {
		t.Errorf("Time = %v after two 0.1 steps", g.Time())
	}
}

func TestScriptFuncMovesAnObject(t *testing.T) {
	g := New()
	o := g.Add(NewObject("mover"))
	o.AddScript(ScriptFunc(func(o *Object, dt float64) {
		o.Transform.Translate(vector.Vector3{X: 10 * dt})
	}))
	for range 10 {
		g.Step(0.1)
	}
	if got := o.Transform.Position.X; got < 9.999 || got > 10.001 {
		t.Errorf("after one second at 10/s the object is at %v", got)
	}
}

func TestTouchingObjectsAreBothTold(t *testing.T) {
	var log []string
	g := New()
	a := g.Add(NewObject("a"))
	a.Collider = collider.Box{X: -5, Y: -5, Width: 10, Height: 10}
	a.AddScript(recorder{&log})
	b := g.Add(NewObject("b"))
	b.Collider = collider.Circle{Radius: 3}
	b.Transform.Position = vector.Vector3{X: 6}
	b.AddScript(recorder{&log})
	far := g.Add(NewObject("far"))
	far.Collider = collider.Circle{Radius: 1}
	far.Transform.Position = vector.Vector3{X: 100}
	g.Step(0)
	hits := 0
	for _, l := range log {
		if l == "a hit b" || l == "b hit a" {
			hits++
		}
		if l == "a hit far" || l == "far hit a" {
			t.Error("objects far apart were told they touch")
		}
	}
	if hits != 2 {
		t.Errorf("log = %v, want a and b told once each", log)
	}
}

func TestDestroyedObjectsLeave(t *testing.T) {
	g := New()
	o := g.Add(NewObject("doomed"))
	o.AddScript(ScriptFunc(func(o *Object, dt float64) { o.Destroy() }))
	g.Add(NewObject("stays"))
	g.Step(0.1)
	if g.Find("doomed") != nil || len(g.Objects()) != 1 || o.Game() != nil {
		t.Error("a destroyed object is still in the world")
	}
}

func TestObjectsAddedDuringAStepJoinTheNext(t *testing.T) {
	g := New()
	updates := 0
	g.Add(NewObject("spawner")).AddScript(ScriptFunc(func(o *Object, dt float64) {
		if g.Find("child") == nil {
			o.Game().Add(NewObject("child")).AddScript(ScriptFunc(func(*Object, float64) { updates++ }))
		}
	}))
	g.Step(0.1)
	if updates != 0 {
		t.Error("an object added mid-step was updated in that step")
	}
	g.Step(0.1)
	if updates != 1 {
		t.Errorf("the added object updated %d times in the next step, want 1", updates)
	}
}

func TestAnimatorFrameIsShownOverTheSprite(t *testing.T) {
	still, _ := sprite.Rect(1, 1, canvas.RGB(1, 1, 1))
	moving, _ := sprite.Rect(1, 1, canvas.RGB(2, 2, 2))
	o := NewObject("x")
	o.Sprite = still
	if o.Look() != still {
		t.Error("an object without an animator does not show its sprite")
	}
	o.Animator = animation.NewAnimator()
	o.Animator.Add("idle", &animation.Animation{Frames: []*sprite.Sprite{moving}, FPS: 1, Loop: true})
	if o.Look() != still {
		t.Error("an animator playing nothing hides the sprite")
	}
	o.Animator.Play("idle")
	if o.Look() != moving {
		t.Error("a playing animation is not what the object shows")
	}
}

func TestDrawOrdersByZ(t *testing.T) {
	red, blue := canvas.RGB(255, 0, 0), canvas.RGB(0, 0, 255)
	g := New()
	top, _ := sprite.Rect(4, 4, red)
	bottom, _ := sprite.Rect(4, 4, blue)
	a := g.Add(NewObject("top"))
	a.Sprite, a.Z = top, 1
	b := g.Add(NewObject("bottom"))
	b.Sprite, b.Z = bottom, 0
	for _, o := range g.Objects() {
		o.Transform.Position = vector.Vector3{X: 5, Y: 5}
	}
	cv, err := canvas.NewCanvas(10, 10)
	if err != nil {
		t.Fatal(err)
	}
	g.Background = canvas.RGB(0, 255, 0)
	g.Draw(cv)
	if got := cv.At(5, 5); got != red {
		t.Errorf("centre = %08x, want the higher Z (red) on top", got)
	}
	if got := cv.At(0, 0); got != g.Background {
		t.Errorf("corner = %08x, want the background", got)
	}
	a.Hidden = true
	g.Draw(cv)
	if got := cv.At(5, 5); got != blue {
		t.Errorf("with the top hidden the centre is %08x, want blue", got)
	}
}

// marker paints one pixel where its object is.
type marker struct{ c canvas.Color }

func (marker) Update(*Object, float64) {}
func (m marker) Paint(o *Object, cv *canvas.Canvas) {
	cv.Pixel(int(o.Transform.Position.X), int(o.Transform.Position.Y), m.c)
}

func TestPaintersDrawAfterTheirSprite(t *testing.T) {
	g := New()
	s, _ := sprite.Rect(4, 4, canvas.RGB(0, 0, 255))
	o := g.Add(NewObject("hud"))
	o.Sprite = s
	o.Transform.Position = vector.Vector3{X: 5, Y: 5}
	white := canvas.RGB(255, 255, 255)
	o.AddScript(marker{white})
	cv, _ := canvas.NewCanvas(10, 10)
	g.Draw(cv)
	if got := cv.At(5, 5); got != white {
		t.Errorf("the painted pixel is %08x, want it over the sprite", got)
	}
}
