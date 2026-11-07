package bakes

import (
	"fmt"
	"math"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
)

// BedModel is the asset a villager sleeps in, the one the houses' beds are.
const BedModel = "vintage_day_bed"

// A sleeper is baked in as few pictures as breathe: SleepFrames of a slow
// breath, in the one way the houses' beds lie, so the bed and the sleeper
// are one picture and the bed hides the body where it should.
const SleepFrames = 4

// Where the bed lies in the sleeper's space, in metres below the floor
// their back is level with, turned so its length runs along the body, and
// where on it they lie.
var (
	bedAt   = [3]float64{-bedAside, 0, 0}
	bedTurn = math.Pi / 2
)

// bedAside is how many metres the bed's middle is from the one who sleeps
// in it, across the bed, so they lie in the front of it by the pillow.
const bedAside = 0.2

// SleepShift is how many pixels, across and down, the bed's middle is
// from the pivot of a sleeper's picture, which stands where the person
// does: the bed's middle, where the game stands the sleeper, is the pivot
// moved by it.
func SleepShift() (dx, dy int) {
	return int(math.Round(bedAside * kit.Cell)), -int(math.Round(bedAside * kit.Cell / 2))
}

// sleepYaw is the turn the sleeper is baked at: the bed's own heading 0,
// 45 degrees, and the quarter turn of the person that lays them along it.
const sleepYaw = 315

// SleepName is the name of frame k of villager id asleep at stage:
// person_<id>_<stage>_sleep_f<frame>.
func SleepName(id string, stage life.Stage, k int) string {
	return SleepMotionName(id, stage, human.Sleep, k)
}

// SleepMotionName is the name of frame k of villager id asleep at stage
// in the sleeping motion mo: in a bed of their own, or on one side of the
// double bed.
func SleepMotionName(id string, stage life.Stage, mo human.Motion, k int) string {
	return fmt.Sprintf("%s_f%d", PersonName(id, stage, mo.Name), k)
}

// isSleep reports whether mo is a way of sleeping, baked in SleepFrames
// pictures and no facing of its own.
func isSleep(mo human.Motion) bool {
	return mo.Name == human.Sleep.Name || mo.Name == human.SleepLeft.Name || mo.Name == human.SleepRight.Name
}

// hasSleep reports whether grown g sleeps in a bed: the grown-up and the
// elder villagers; the player has a bed of their own to use.
func hasSleep(g kit.Grown) bool { return hasMotion(g.Motions, human.Sleep) }

// sleepTask bakes grown g of villager v asleep in the way mo says:
// SleepFrames pictures, one task.
func (b *Baker) sleepTask(v kit.Villager, g kit.Grown, mo human.Motion) Task {
	name := fmt.Sprintf("person %s %s %s %s", versionOf(mo), v.ID, g.Stage, mo.Name)
	return Task{Name: name, Kind: "person", ID: v.ID, Run: func() error {
		body, err := b.bodyOf(v, g)
		if err != nil {
			return fmt.Errorf("%s as a %s: %w", v.ID, g.Stage, err)
		}
		return b.each(SleepFrames, func(k int) error {
			shot := bake.Shot{Animation: mo.Name, Time: mo.Length * float64(k) / SleepFrames, Yaw: sleepYaw}
			return b.Shelf.Put(SleepMotionName(v.ID, g.Stage, mo, k), sleeper(body.model, mo, shot))
		})
	}}
}

// sleeper is the picture of the sleeper as shot: with the bed they lie in,
// but on one side of a double bed, without it.
func sleeper(m *gltf.Model, mo human.Motion, sh bake.Shot) bake.Sprite {
	sp := shoot(m, frame{shot: sh})
	if mo.Name == human.Sleep.Name {
		return sp
	}
	return withoutBed(sp, bareBed(m, mo, sh))
}

// layInBed gives m, a person already animated, the bed to sleep in, if
// they sleep.
func (b *Baker) layInBed(m *gltf.Model, motions []human.Motion) error {
	if err := b.layInDoubleBed(m, motions); err != nil || !hasMotion(motions, human.Sleep) {
		return err
	}
	var asset kit.Asset
	for _, a := range b.Manifest.Assets {
		if a.ID == BedModel {
			asset = a
		}
	}
	if asset.ID == "" {
		return fmt.Errorf("no asset %s to sleep in", BedModel)
	}
	bed, err := Load(asset)
	if err != nil {
		return fmt.Errorf("%s: %w", asset.ID, err)
	}
	return human.Lay(m, bed, human.BedPlace{At: bedAt, Yaw: bedTurn, Scale: scaleOf(asset), During: human.Sleep.Name})
}

// scaleOf is how many times its real size asset a is drawn.
func scaleOf(a kit.Asset) float64 {
	if a.Scale == 0 {
		return 1
	}
	return a.Scale
}
