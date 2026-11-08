package game

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// lazyObject is an object drawn with a 4 by 4 sprite made on first use,
// counting the makes.
func lazyObject(g *Game, x float64, made *int) *Object {
	o := NewObject("lazy")
	o.Sprite = &sprite.Sprite{Make: func() *canvas.Canvas {
		*made++
		cv, _ := canvas.NewCanvas(4, 4)
		cv.Clear(canvas.RGB(200, 30, 30))
		return cv
	}}
	o.Transform.Position = vector.Vector3{X: x, Y: 10}
	g.Add(o)
	return o
}

// A sprite made on first use, far outside the canvas, is not made until
// it comes near, and the picture is the same as when every sprite is made.
func TestAFarSpriteIsNotMadeUntilItComesNear(t *testing.T) {
	frame := func(reach float64, x float64) ([]canvas.Color, int) {
		g := New()
		g.LazyReach = reach
		made := 0
		lazyObject(g, 5, &made)
		far := lazyObject(g, 2000, &made)
		cv, _ := canvas.NewCanvas(64, 64)
		g.Draw(cv)
		if reach > 0 && made != 1 {
			t.Fatalf("made %d sprites with one far away", made)
		}
		far.Transform.Position.X = x
		g.Draw(cv)
		return slices.Clone(cv.Pixels), made
	}
	with, madeWith := frame(100, 30)
	without, madeWithout := frame(0, 30)
	if !slices.Equal(with, without) {
		t.Fatal("the picture differs when far sprites are not made")
	}
	if madeWith != 2 || madeWithout != 2 {
		t.Fatalf("made %d and %d after it came near", madeWith, madeWithout)
	}
}
