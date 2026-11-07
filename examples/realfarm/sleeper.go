package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/sprite"
)

// SleepBaked says the sleeper is baked with the bed, by bakes.layInBed.
func (a *realArt) SleepBaked() bool { return true }

// asleep is the clip of a villager lying in a bed, whichever way they
// face: the frames of the one way the houses' beds lie, each stood so
// that the bed's middle is where the villager is put. A villager whose
// sleep is not baked yet lies in no clip; a blank one stands in for it.
func (b *realBody) asleep(mo human.Motion) farm.Clip {
	if !b.sleepsIn(mo) {
		return farm.Clip{Frames: []*sprite.Sprite{{Image: newFrame()}}, FPS: 1}
	}
	dx, dy := bakes.SleepShift()
	if mo.Name != human.Sleep.Name {
		dx, dy = bakes.SleepShiftSide(mo)
	}
	out := make([]*sprite.Sprite, bakes.SleepFrames)
	for k := range out {
		name := bakes.SleepMotionName(b.id, b.stage, mo, k)
		at, _ := b.art.shelf.Sprite(name)
		out[k] = &sprite.Sprite{PivotX: float64(at.PivotX+dx) / float64(at.W), PivotY: float64(at.PivotY+dy) / float64(at.H),
			Make: func() *canvas.Canvas { return b.art.picture(name) }}
	}
	return farm.Clip{Frames: out, FPS: bakes.SleepFrames / human.Sleep.Length}
}

// Sleeps is whether the frames of the villager asleep are baked.
func (b *realBody) Sleeps() bool { return b.sleepsIn(human.Sleep) }

// SleepsAs is whether the frames of the villager asleep on the side of the
// double bed that anim says are baked.
func (b *realBody) SleepsAs(anim string) bool { return b.sleepsIn(sleepMotion(anim)) }

// sleepsIn is whether the frames of the villager asleep in the way mo
// says are baked.
func (b *realBody) sleepsIn(mo human.Motion) bool {
	return b.art.shelf.Has(bakes.SleepMotionName(b.id, b.stage, mo, bakes.SleepFrames-1))
}

// sleepMotion is the way of sleeping the animation anim is.
func sleepMotion(anim string) human.Motion {
	switch anim {
	case farm.AnimSleepLeft:
		return human.SleepLeft
	case farm.AnimSleepRight:
		return human.SleepRight
	}
	return human.Sleep
}
