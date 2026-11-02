package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// savePNG writes cv to path.
func savePNG(t *testing.T, cv *canvas.Canvas, path string) {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, cv.Width, cv.Height))
	for y := range cv.Height {
		for x := range cv.Width {
			c := cv.At(x, y)
			img.SetNRGBA(x, y, color.NRGBA{c.R(), c.G(), c.B(), 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

// With SGE_SHOT set to a folder, paints the real creation screen, the real
// person on the stage, in both framings to PNGs there, and logs how long
// each step takes. Not run otherwise: it needs the human data in the
// cache, and seconds.
func TestRealCreatorScreens(t *testing.T) {
	dir := os.Getenv("SGE_SHOT")
	if dir == "" {
		t.Skip("SGE_SHOT not set")
	}
	b, err := newBaker()
	if err != nil {
		t.Fatal(err)
	}
	live := newLivePreview(bakerMaker{b}, 0)
	defer live.Close()
	c := newCreator(testArt(t), live, func(any) {})
	if os.Getenv("SGE_SOFT") != "" {
		c.stage = newLiveStage(softDevice)
	}
	defer c.stage.Close()
	cv, _ := canvas.NewCanvas(1600, 900)
	c.paint(cv)
	t0 := time.Now()
	for live.Busy() && time.Since(t0) < 90*time.Second {
		time.Sleep(20 * time.Millisecond)
	}
	t.Logf("first person, kit and textures read: %v (err %v)", time.Since(t0), live.Err())
	t1 := time.Now()
	c.advance(1.0 / 60)
	t.Logf("upload: %v, device %T", time.Since(t1), c.stage.dev)
	frame := func(name string) {
		for range 120 {
			c.advance(1.0 / 60)
		}
		t2 := time.Now()
		c.paint(cv)
		t.Logf("%s: first paint (render) %v", name, time.Since(t2))
		savePNG(t, cv, filepath.Join(dir, name+".png"))
	}
	c.pick(0)
	c.row = 0
	c.handle(cin{Select: true})
	c.handle(cin{Typed: "Ana"})
	frame("body")
	c.pick(indexOfCategory(t, c, "Cabelo"))
	c.refresh()
	frame("head")
	// colours: time one from the keypress to the picture
	t3 := time.Now()
	c.row = 1
	c.handle(cin{Dx: 3})
	waitFor(t, "recolour planned", func() bool { return !live.Busy() })
	t.Logf("hair colour planned: %v", time.Since(t3))
	t4 := time.Now()
	c.advance(1.0 / 60)
	c.paint(cv)
	t.Logf("hair colour applied and drawn: %v", time.Since(t4))
	c.pick(indexOfCategory(t, c, "Pele"))
	c.row = 3
	c.handle(cin{Dx: 4})
	waitFor(t, "skin", func() bool { return !live.Busy() })
	c.advance(1.0 / 60)
	c.pick(indexOfCategory(t, c, "Roupa"))
	c.row = 0
	t5 := time.Now()
	c.handle(cin{Dx: 1})
	waitFor(t, "remake", func() bool { return !live.Busy() })
	t.Logf("another garment, made again: %v", time.Since(t5))
	c.advance(1.0 / 60)
	c.view.Turn(35)
	frame("body-turned")
}
