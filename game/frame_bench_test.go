package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// benchWorld is a busy frame: n spinning, moving objects with sprites and
// colliders, each told what it bumps into, the way a small game has them.
func benchWorld(b *testing.B, n int) *Game {
	g := New()
	look, err := sprite.Rect(12, 12, canvas.RGB(250, 200, 40))
	if err != nil {
		b.Fatal(err)
	}
	for i := range n {
		o := g.Add(NewObject("thing"))
		o.Sprite, o.Collider = look, collider.Circle{Radius: 6}
		o.Transform.Position = vector.Vector3{X: float64(20 + (i*37)%440), Y: float64(20 + (i*53)%280)}
		o.AddScript(ScriptFunc(func(o *Object, dt float64) {
			o.Transform.Rotate(vector.Vector3{Z: 1}, 2*dt)
			o.Transform.Translate(vector.Vector3{X: 10 * dt})
		}))
		o.AddScript(&bumps{})
	}
	return g
}

// BenchmarkStep is the world's update of one frame: scripts and collisions.
func BenchmarkStep(b *testing.B) {
	g := benchWorld(b, 200)
	for b.Loop() {
		g.Step(1.0 / 60)
	}
}

// BenchmarkDraw is the world painted once onto a window-sized canvas.
func BenchmarkDraw(b *testing.B) {
	g := benchWorld(b, 200)
	cv, err := canvas.NewCanvas(480, 320)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		g.Draw(cv)
	}
}

// bumps counts what its object touches.
type bumps struct{ n int }

func (b *bumps) Update(*Object, float64)    {}
func (b *bumps) Collision(o, other *Object) { b.n++ }
