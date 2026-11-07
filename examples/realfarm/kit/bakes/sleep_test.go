package bakes

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// Only the children, teens, grown-up and elder villagers sleep in a bed, never the player,
// in one task of few frames each, and no facing of its own: the cost that
// every villager's sleep adds to the bake is those frames.
func TestOnlyGrownVillagersSleepAndInOneTask(t *testing.T) {
	var b Baker
	sleeps, frames := 0, 0
	for _, task := range b.peopleTasks() {
		if !strings.HasSuffix(task.Name, " sleep") {
			if strings.Contains(task.Name, " sleep ") {
				t.Errorf("a sleep baked by facing: %s", task.Name)
			}
			continue
		}
		sleeps++
		frames += SleepFrames
		if strings.Contains(task.Name, " player ") || task.Core || task.Kind != "person" {
			t.Errorf("task %+v", task)
		}
	}
	if want := 4 * (len(kit.Villagers) - 1); sleeps != want {
		t.Errorf("%d sleep tasks, not %d: the children, teens, grown and elders of all but the player", sleeps, want)
	}
	if SleepFrames > 4 || frames > SleepFrames*4*(len(kit.Villagers)-1) {
		t.Errorf("%d frames a sleeper, %d in all", SleepFrames, frames)
	}
}

// Sleeping has a version of its own, apart from every other motion's, and
// leaves the motions' version as it was.
func TestSleepHasItsOwnVersion(t *testing.T) {
	if versionOf(human.Sleep) == versionOf(human.Stand) || versionOf(human.Sleep) == versionOf(human.Sit) {
		t.Errorf("sleep's version %s is another's", versionOf(human.Sleep))
	}
	if versionOf(human.Stand) != "v4" {
		t.Errorf("the motions' version changed: %s", versionOf(human.Stand))
	}
}

// The bed's middle is off the sleeper's pivot by what its lying aside
// makes it, as measured on the baked pictures: 6 across, 3 up.
func TestTheBedIsOffTheSleepersPivot(t *testing.T) {
	if dx, dy := SleepShift(); dx != 6 || dy != -3 {
		t.Errorf("shift %d, %d", dx, dy)
	}
	if SleepName("ana", life.Adult, 2) != "person_ana_adult_sleep_f2" {
		t.Errorf("name %s", SleepName("ana", life.Adult, 2))
	}
}

// A couple's two sides of the double bed are each a task of their own, one
// way and a few frames, for the grown and the elders of all but the player.
func TestACoupleSleepsInTheDoubleBedOnEachSide(t *testing.T) {
	var b Baker
	sides := map[string]int{}
	for _, task := range b.peopleTasks() {
		for _, mo := range human.SleepingTogether {
			if strings.HasSuffix(task.Name, " "+mo.Name) {
				sides[mo.Name]++
			}
		}
	}
	for _, mo := range human.SleepingTogether {
		if want := 2 * (len(kit.Villagers) - 1); sides[mo.Name] != want {
			t.Errorf("%d tasks of %s, not %d", sides[mo.Name], mo.Name, want)
		}
	}
	if versionOf(human.SleepLeft) != versionOf(human.Sleep) || SleepMotionName("ana", life.Adult, human.SleepRight, 1) != "person_ana_adult_sleep_r_f1" {
		t.Errorf("version %s, name %s", versionOf(human.SleepLeft), SleepMotionName("ana", life.Adult, human.SleepRight, 1))
	}
}

// What a sleeper's picture has that the bare bed has not is the sleeper:
// the pixels the same are taken away, and the picture is cut down to the rest.
func TestWithoutTheBedKeepsOnlyTheSleeper(t *testing.T) {
	sp := bake.Sprite{Image: image.NewNRGBA(image.Rect(0, 0, 4, 4)), Pivot: image.Pt(2, 2)}
	bare := bake.Sprite{Image: image.NewNRGBA(image.Rect(0, 0, 4, 4)), Pivot: image.Pt(2, 2)}
	bed, body := color.NRGBA{100, 50, 0, 255}, color.NRGBA{0, 0, 200, 255}
	for y := range 4 {
		for x := range 4 {
			sp.Image.SetNRGBA(x, y, bed)
			bare.Image.SetNRGBA(x, y, bed)
		}
	}
	sp.Image.SetNRGBA(1, 2, body)
	sp.Image.SetNRGBA(2, 2, body)
	got := withoutBed(sp, bare)
	if got.Image.Bounds().Dx() != 2 || got.Image.Bounds().Dy() != 1 || got.Pivot != image.Pt(1, 0) || got.Image.NRGBAAt(0, 0) != body {
		t.Errorf("%v pivot %v", got.Image.Bounds(), got.Pivot)
	}
}
