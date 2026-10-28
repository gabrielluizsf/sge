package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// The camera moves the world, not what is fixed to the screen, and a
// painter finds where its object is on screen with ToScreen.
func TestTheCameraMovesTheWorldButNotTheScreen(t *testing.T) {
	red, blue, white := canvas.RGB(255, 0, 0), canvas.RGB(0, 0, 255), canvas.RGB(255, 255, 255)
	g := New()
	look, _ := sprite.Rect(2, 2, red)
	thing := g.Add(NewObject("thing"))
	thing.Sprite = look
	thing.Transform.Position = vector.Vector3{X: 105, Y: 53}
	thing.AddScript(screenMarker{white})
	hudLook, _ := sprite.Rect(2, 2, blue)
	hud := g.Add(NewObject("hud"))
	hud.Sprite, hud.Screen = hudLook, true
	hud.Transform.Position = vector.Vector3{X: 2, Y: 2}
	g.Camera = vector.Vector3{X: 100, Y: 50}
	cv, _ := canvas.NewCanvas(10, 10)
	g.Draw(cv)
	if got := cv.At(4, 2); got != red {
		t.Errorf("the thing at 105,53 seen from 100,50 is not at 4..5,2..3: %08x", got)
	}
	if got := cv.At(5, 3); got != white {
		t.Errorf("the painter did not paint at its object's place on screen: %08x", got)
	}
	if got := cv.At(1, 1); got != blue {
		t.Errorf("the HUD moved with the camera: %08x", got)
	}
	if w, h := g.View(); w != 10 || h != 10 {
		t.Errorf("view %dx%d, not the canvas drawn on", w, h)
	}
}

// Seen from a camera between pixels, every object is moved by the same
// whole number of pixels, so things side by side do not shimmer apart.
func TestTheCameraMovesEverythingByWholePixels(t *testing.T) {
	g := New()
	red := canvas.RGB(255, 0, 0)
	for _, x := range []float64{3, 5} {
		s, _ := sprite.Rect(2, 2, red)
		s.PivotX, s.PivotY = 0, 0
		o := g.Add(NewObject("tile"))
		o.Sprite = s
		o.Transform.Position = vector.Vector3{X: x}
	}
	cv, _ := canvas.NewCanvas(10, 2)
	for _, cam := range []float64{0.4, 0.5, 0.6, 1.49} {
		g.Camera = vector.Vector3{X: cam}
		g.Draw(cv)
		for x := range 10 {
			if cv.At(x, 0) != red && x >= 3-int(cam+0.5) && x < 7-int(cam+0.5) {
				t.Errorf("camera at %v: a gap at %d", cam, x)
			}
		}
	}
}

// screenMarker paints a pixel where its object is on screen.
type screenMarker struct{ c canvas.Color }

func (screenMarker) Update(*Object, float64) {}
func (m screenMarker) Paint(o *Object, cv *canvas.Canvas) {
	p := o.Game().ToScreen(o.Transform.Position)
	cv.Pixel(int(p.X), int(p.Y), m.c)
}
