package main

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
)

// liveCreator is a creation screen over a fake maker, drawn on the
// software device and painted once, so it knows where things are.
func liveCreator(t *testing.T, f *fakeMaker) *creator {
	t.Helper()
	live := newLivePreview(f, 0)
	t.Cleanup(live.Close)
	c := newCreator(testArt(t), live, func(any) {})
	c.stage = newLiveStage(softDevice)
	cv, _ := canvas.NewCanvas(1280, 720)
	c.paint(cv)
	return c
}

// idle is a frame in which the player does nothing.
func (c *creator) idle() { c.handle(cin{}); c.advance(1.0 / 60) }

func TestDraggingTheStageTurnsThePerson(t *testing.T) {
	c := liveCreator(t, &fakeMaker{})
	x, y := c.stageAt.X+c.stageAt.W/2, c.stageAt.Y+c.stageAt.H/2
	c.handle(cin{X: x, Y: y, Press: true, Held: true})
	c.handle(cin{X: x + 40, Y: y, Held: true})
	c.handle(cin{X: x + 60, Y: y + 30, Held: true})
	if want := 60 * turnPerPixel; c.view.yaw != want {
		t.Errorf("a 60 pixel drag turned %v degrees, want %v", c.view.yaw, want)
	}
	c.handle(cin{X: x + 60, Y: y})
	c.handle(cin{X: x - 200, Y: y, Held: true}) // a new press is needed to turn again
	if c.view.yaw != 60*turnPerPixel {
		t.Errorf("turned after the mouse was let go: %v", c.view.yaw)
	}
}

func TestNothingElseTurnsThePerson(t *testing.T) {
	f := &fakeMaker{}
	c := liveCreator(t, f)
	for range 120 {
		c.idle()
	}
	c.handle(cin{Dx: 1})
	c.handle(cin{Dy: 1})
	c.handle(cin{Select: true})
	c.handle(cin{Press: true, Held: true, X: 5, Y: 5}) // a press off the stage
	c.handle(cin{Held: true, X: 300, Y: 5})
	time.Sleep(50 * time.Millisecond)
	for range 30 {
		c.idle()
	}
	if c.view.yaw != 0 {
		t.Errorf("the person turned %v degrees on its own", c.view.yaw)
	}
}

func TestADragThatStartsOnARowDoesNotTurnThePerson(t *testing.T) {
	c := liveCreator(t, &fakeMaker{})
	var on hit
	for _, h := range c.hits {
		if h.kind == hitRow {
			on = h
		}
	}
	c.handle(cin{X: on.at.X + 3, Y: on.at.Y + 3, Press: true, Held: true})
	c.handle(cin{X: c.stageAt.X + 100, Y: on.at.Y, Held: true})
	if c.view.yaw != 0 {
		t.Errorf("turned %v", c.view.yaw)
	}
}

func TestKeysTurnThePerson(t *testing.T) {
	c := liveCreator(t, &fakeMaker{})
	c.dt = 0.5
	c.handle(cin{Turn: 1})
	if c.view.yaw != keyTurn*0.5 {
		t.Errorf("yaw %v", c.view.yaw)
	}
}

// The face categories frame the head, the body and clothes the whole body.
func TestFaceCategoriesFrameTheHead(t *testing.T) {
	c := liveCreator(t, &fakeMaker{})
	for i, cat := range c.cats {
		c.handle(cin{X: -1, Y: -1}) // nothing
		c.pick(i)
		c.refresh()
		want := 0.0
		switch cat.name {
		case "Skin", "Face", "Hair", "Eyes":
			want = 1
		}
		if c.view.goal != want {
			t.Errorf("%s frames %v, want %v", cat.name, c.view.goal, want)
		}
	}
	c.pick(0)
	c.refresh()
	for range 200 {
		c.view.Step(1.0 / 60)
	}
	if c.view.head != 0 {
		t.Errorf("the body category left the camera at %v", c.view.head)
	}
}

// Changing a colour shows at once and does not make the mesh again;
// changing what is worn does.
func TestColourRowsDoNotRemakeTheMesh(t *testing.T) {
	f := &fakeMaker{}
	c := liveCreator(t, f)
	waitFor(t, "first person", func() bool { return !c.live.Busy() })
	c.advance(1.0 / 60)
	hair := indexOfCategory(t, c, "Hair")
	c.pick(hair)
	c.row = 1 // the colour
	for range 6 {
		c.handle(cin{Dx: 1})
	}
	skin := indexOfCategory(t, c, "Skin")
	c.pick(skin)
	c.row = 3 // the tone
	c.handle(cin{Dx: 1})
	c.handle(cin{Dx: 1})
	eyes := indexOfCategory(t, c, "Eyes")
	c.pick(eyes)
	c.handle(cin{Dx: 1})
	waitFor(t, "recoloured", func() bool { return !c.live.Busy() })
	if f.makes() != 1 {
		t.Fatalf("the mesh was made %d times for colours", f.makes())
	}
	c.pick(indexOfCategory(t, c, "Hair"))
	c.row = 0 // the style
	c.handle(cin{Dx: 1})
	waitFor(t, "remade", func() bool { return f.makes() == 2 && !c.live.Busy() })
}

func indexOfCategory(t *testing.T, c *creator, name string) int {
	t.Helper()
	for i, cat := range c.cats {
		if cat.name == name {
			return i
		}
	}
	t.Fatalf("no category %s", name)
	return 0
}

// Input and the frame go on while a person is being made.
func TestTheScreenDoesNotWaitForABuild(t *testing.T) {
	f := &fakeMaker{gate: make(chan struct{})}
	defer close(f.gate)
	c := liveCreator(t, f)
	start := time.Now()
	for range 50 {
		c.handle(cin{Dx: 1})
		c.advance(1.0 / 60)
	}
	if d := time.Since(start); d > 200*time.Millisecond {
		t.Errorf("50 frames took %v", d)
	}
	if !c.live.Busy() {
		t.Error("the build was not under way")
	}
}
