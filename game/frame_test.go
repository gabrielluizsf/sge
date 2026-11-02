package game

import (
	"slices"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sound"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// dot is a painter that draws a dot, and has a version if it says so.
type dot struct {
	version *uint64
	drawn   *int
}

func (dot) Update(*Object, float64) {}
func (d dot) Paint(_ *Object, cv *canvas.Canvas) {
	*d.drawn++
	cv.Put(1, 1, canvas.RGB(255, 0, 0))
}

// bare is a painter that draws a dot and has no version.
type bare struct{ drawn *int }

func (bare) Update(*Object, float64) {}
func (b bare) Paint(_ *Object, cv *canvas.Canvas) {
	*b.drawn++
	cv.Put(1, 1, canvas.RGB(255, 0, 0))
}

// Version is the dot's.
func (d dot) Version() uint64 { return *d.version }

// framed is a game with a sprite that stays, and the canvas it is drawn on.
func framed(t *testing.T) (*Game, *Object, *canvas.Canvas) {
	g := New()
	g.Background = canvas.RGB(10, 20, 30)
	o := g.Add(NewObject("thing"))
	o.Sprite, _ = sprite.Rect(4, 4, canvas.RGB(200, 100, 50))
	o.Transform.Position = vector.Vector3{X: 8, Y: 8}
	cv, err := canvas.NewCanvas(32, 24)
	if err != nil {
		t.Fatal(err)
	}
	return g, o, cv
}

// A frame that would come out the pixels on the canvas is not drawn, and
// any change to what is seen is.
func TestAFrameTheSameIsNotDrawn(t *testing.T) {
	g, o, cv := framed(t)
	if !g.Frame(cv, 0.1) {
		t.Fatal("the first frame was not drawn")
	}
	for range 5 {
		if g.Frame(cv, 0.1) {
			t.Fatal("a frame of a world that stands still was drawn")
		}
	}
	// In order: each change is seen from where the one before left the world.
	for _, c := range []struct {
		name   string
		change func()
	}{
		{"the object moves", func() { o.Transform.Position.X++ }},
		{"the object hides", func() { o.Hidden = true }},
		{"it shows again", func() { o.Hidden = false }},
		{"the camera moves", func() { g.Camera.Y += 2 }},
		{"the background", func() { g.Background = canvas.RGB(1, 2, 3) }},
		{"another object", func() { g.Add(NewObject("dot")).Sprite = o.Sprite }},
		{"the sprite changes", func() { o.Sprite, _ = sprite.Rect(4, 4, canvas.RGB(1, 1, 1)) }},
		{"an object is gone", func() { g.Objects()[1].Destroy() }},
		{"invalidated", g.Invalidate},
	} {
		c.change()
		if !g.Frame(cv, 0.1) {
			t.Errorf("%s: the frame was not drawn", c.name)
		}
		if g.Frame(cv, 0.1) {
			t.Errorf("%s: the frame after it was drawn", c.name)
		}
	}
	other, _ := canvas.NewCanvas(32, 24)
	if !g.Frame(other, 0.1) {
		t.Error("another canvas was not drawn")
	}
}

// A frame left undrawn leaves the canvas as a full draw would make it.
func TestAFrameNotDrawnIsTheFrameDrawn(t *testing.T) {
	g, o, cv := framed(t)
	scratch, _ := canvas.NewCanvas(32, 24)
	for i := range 20 {
		if i == 7 {
			o.Transform.Position.X += 3
		}
		g.Frame(cv, 0.1)
		g.Draw(scratch)
		if !slices.Equal(cv.Pixels, scratch.Pixels) {
			t.Fatalf("frame %d is not what drawing it whole makes", i)
		}
	}
}

// A painter that is not a Versioned is drawn every frame, one that is
// only when its version moves on.
func TestPaintersAreDrawnUnlessVersioned(t *testing.T) {
	g, _, cv := framed(t)
	drawn := 0
	g.Add(NewObject("bare")).AddScript(bare{&drawn})
	for range 4 {
		g.Frame(cv, 0.1)
	}
	if drawn != 4 {
		t.Errorf("a painter with no version was drawn %d times in 4 frames", drawn)
	}
	g2, _, cv2 := framed(t)
	version, drawn := uint64(1), 0
	g2.Add(NewObject("versioned")).AddScript(dot{&version, &drawn})
	for range 4 {
		g2.Frame(cv2, 0.1)
	}
	if drawn != 1 {
		t.Errorf("a painter whose version stood was drawn %d times in 4 frames", drawn)
	}
	version++
	g2.Frame(cv2, 0.1)
	if drawn != 2 {
		t.Errorf("a painter that moved on was not drawn: %d", drawn)
	}
}

// Frames left undrawn still step the world, and the sounds its scripts
// start still play: the mixer is fed whatever is drawn.
func TestFramesNotDrawnStillStepAndSound(t *testing.T) {
	g, _, cv := framed(t)
	g.Sound = sound.NewMixer(sound.Rate)
	g.Sound.SetRate(sound.Rate)
	beep := &sound.Clip{Rate: sound.Rate, Samples: sound.Tone{From: 440, Seconds: 0.2}.Make()}
	steps, drawn := 0, 0
	g.Add(NewObject("ear")).AddScript(ScriptFunc(func(_ *Object, _ float64) {
		steps++
		if steps%10 == 0 {
			g.Sound.Play(beep, sound.Voice{Gain: 1})
		}
	}))
	heard := 0
	block := make([]float32, 2*735)
	for range 60 {
		if g.Frame(cv, 1.0/60) {
			drawn++
		}
		g.Sound.Fill(block)
		if slices.ContainsFunc(block, func(s float32) bool { return s != 0 }) {
			heard++
		}
	}
	if steps != 60 {
		t.Errorf("the world stepped %d times in 60 frames", steps)
	}
	if drawn != 1 {
		t.Errorf("a world that stands still was drawn %d times", drawn)
	}
	if heard < 50 {
		t.Errorf("only %d of 60 frames had sound to play: the mixer was not fed", heard)
	}
}
