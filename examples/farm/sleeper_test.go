package main

import (
	"image"
	"testing"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
	"github.com/gabrielluizsf/sge/farm"
	"github.com/gabrielluizsf/sge/life"
)

// The realistic look draws the bed with its sleeper, and a villager's
// sleep shows only once it is baked: until then there is no sleeper, so
// the bed is not taken away for nothing. Once it is, the frames are the
// bed's, each stood so that the bed's middle is where the villager is.
func TestASleeperIsDrawnOnceBaked(t *testing.T) {
	shelf, err := bakes.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	a := newArt(shelf)
	if s, ok := a.(farm.SleepArt); !ok || !s.SleepBaked() {
		t.Fatal("the bed is not drawn with the sleeper")
	}
	who := farm.Who{Woman: true, Role: "farmer", Stage: life.Adult}
	body := a.Person(who)
	if sl, ok := body.(farm.Sleeper); !ok || sl.Sleeps() {
		t.Fatal("asleep before it is baked")
	}
	id := villagerOf(who)
	for k := range bakes.SleepFrames {
		if err := shelf.Put(bakes.SleepName(id, who.Stage, k), bake.Sprite{Image: image.NewNRGBA(image.Rect(0, 0, 100, 50)), Pivot: image.Pt(40, 30)}); err != nil {
			t.Fatal(err)
		}
	}
	body = a.Person(who)
	if !body.(farm.Sleeper).Sleeps() {
		t.Fatal("not asleep once baked")
	}
	c := body.Clip(farm.AnimSleep, "left")
	dx, dy := bakes.SleepShift()
	if len(c.Frames) != bakes.SleepFrames || c.FPS <= 0 {
		t.Fatalf("%d frames at %v", len(c.Frames), c.FPS)
	}
	if got, want := c.Frames[0].PivotX, float64(40+dx)/100; got != want {
		t.Errorf("pivot across %v, not %v", got, want)
	}
	if got, want := c.Frames[0].PivotY, float64(30+dy)/50; got != want {
		t.Errorf("pivot down %v, not %v", got, want)
	}
}
