package game

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// An object that Fades is drawn faded while it covers a Reveal behind it,
// and whole once the Reveal stands in front.
func TestAnObjectFadesOverARevealBehindIt(t *testing.T) {
	draw := func(at vector.Vector3) canvas.Color {
		g := New()
		g.Projection = iso.New(32, 64, 32)
		img, _ := canvas.NewCanvas(40, 80)
		img.FillRect(canvas.Area{Width: 40, Height: 80}, canvas.RGB(255, 0, 0))
		house := NewObject("house")
		house.Sprite = sprite.New(img)
		house.Sprite.PivotX, house.Sprite.PivotY = 0.5, 1
		house.Fades = true
		house.Transform.Position = vector.Vector3{X: 64, Y: 64}
		g.Add(house)
		g.Reveals = func() []Reveal { return []Reveal{{At: at, Height: 40}} }
		cv, _ := canvas.NewCanvas(200, 200)
		g.Camera = g.ToWorld(vector.Vector3{X: -100, Y: -100})
		g.Draw(cv)
		s := g.ToScreen(house.Transform.Position)
		return cv.At(int(s.X), int(s.Y)-20)
	}
	if got := draw(vector.Vector3{X: 32, Y: 32}); got.R() > 160 || got.R() == 0 {
		t.Errorf("house over the farmer is %v, want faded", got)
	}
	if got := draw(vector.Vector3{X: 96, Y: 96}); got != canvas.RGB(255, 0, 0) {
		t.Errorf("house behind the farmer is %v, want whole", got)
	}
}
