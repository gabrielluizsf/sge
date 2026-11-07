package main

import (
	"fmt"
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/sprite"
)

// realBody is a villager at a stage of life, drawn from their baked
// frames, facing eight ways.
type realBody struct {
	art   *realArt
	id    string
	stage life.Stage
}

func (b *realBody) Facings() int { return bakes.Facings }

// motions are the motions that stand in for each animation, the first
// that is baked winning: using what one holds is hoeing, and what is not
// baked yet is done in the nearest that is.
var motions = map[string][]string{
	farm.AnimStand: {"stand"},
	farm.AnimWalk:  {"walk", "stand"},
	farm.AnimRun:   {"run", "walk", "stand"},
	farm.AnimAct:   {"hoe", "chop", "stand"},
	farm.AnimChop:  {"chop", "hoe", "stand"},
	farm.AnimForge: {"forge", "chop", "hoe", "stand"},
	farm.AnimHoe:   {"hoe", "chop", "stand"},
	farm.AnimMine:  {"mine", "chop", "hoe", "stand"},
	farm.AnimWater: {"water", "hoe", "stand"},
	farm.AnimPick:  {"pick", "hoe", "stand"},
	farm.AnimEat:   {"eat", "stand"},
	farm.AnimWell:  {"well", "water", "stand"},
	farm.AnimDip:   {"dip", "water", "stand"},
	farm.AnimSwim:  {"swim", "tread", "walk", "stand"},
	farm.AnimDive:  {"dive", "swim", "stand"},
	farm.AnimRise:  {"rise", "stand"},
	farm.AnimSit:   {"sit", "stand"},
	// Milked sitting, whichever is baked yet.
	farm.AnimMilkDown: {"milkdown", "sit", "stand"},
	farm.AnimMilk:     {"milk", "sit", "stand"},
}

// turnOf is the facing the baker baked that looks the way dir says on
// screen: the baker turns a villager 45 degrees a facing, from the one
// that looks at the right of the screen.
func turnOf(dir string) int {
	return (slices.Index(farm.Facings(8), dir) + turnOffset) % bakes.Facings
}

// Clip is anim facing dir: the frames of the motion that stands for it.
func (b *realBody) Clip(anim, dir string) farm.Clip {
	if mo := sleepMotion(anim); anim == farm.AnimSleep || mo.Name != human.Sleep.Name {
		return b.asleep(mo)
	}
	if anim == farm.AnimMilkUp { // getting up is sitting down, backward
		c := b.Clip(farm.AnimMilkDown, dir)
		slices.Reverse(c.Frames)
		return c
	}
	d := turnOf(dir)
	for _, m := range motions[anim] {
		prefix := bakes.PersonName(b.id, b.stage, m)
		if !b.art.shelf.Has(fmt.Sprintf("%s_d%d_f%d", prefix, d, bakes.Frames-1)) {
			continue
		}
		return farm.Clip{Frames: b.frames(prefix, d), FPS: fpsOf(anim, m)}
	}
	return farm.Clip{Frames: []*sprite.Sprite{{Image: newFrame()}}, FPS: 1}
}

// frames are the sprites of a motion facing d, each read from the shelf
// when it is first drawn, standing on the villager's feet.
func (b *realBody) frames(prefix string, d int) []*sprite.Sprite {
	out := make([]*sprite.Sprite, bakes.Frames)
	for k := range out {
		name := fmt.Sprintf("%s_d%d_f%d", prefix, d, k)
		at, _ := b.art.shelf.Sprite(name)
		out[k] = &sprite.Sprite{PivotX: float64(at.PivotX) / float64(at.W), PivotY: float64(at.PivotY) / float64(at.H),
			Make: func() *canvas.Canvas { return b.art.picture(name) }}
	}
	return out
}

// fpsOf is how many of a motion's frames are shown a second for anim: a
// loop at the pace the motion was made at, a use of what one holds or a
// stroke of the axe once in the time the farm gives it.
func fpsOf(anim, motion string) float64 {
	switch anim {
	case farm.AnimForge:
		return bakes.Frames / farm.ForgeSeconds
	case farm.AnimHoe:
		return bakes.Frames / farm.HoeSeconds
	case farm.AnimWater:
		return bakes.Frames / farm.WaterSeconds
	case farm.AnimPick:
		return bakes.Frames / farm.PickSeconds
	case farm.AnimMine:
		return bakes.Frames / farm.MineSeconds
	case farm.AnimEat:
		return bakes.Frames / farm.EatSeconds
	case farm.AnimWell:
		return bakes.Frames / farm.WellSeconds
	case farm.AnimDip:
		return bakes.Frames / farm.DipSeconds
	case farm.AnimMilkDown:
		return bakes.Frames / farm.MilkDownSeconds
	}
	if anim == farm.AnimAct || anim == farm.AnimChop {
		return bakes.Frames / farm.ActSeconds
	}
	for _, m := range slices.Concat(human.Motions, human.Sitting, human.Mining, human.Eating, human.Drawing, human.Milking) {
		if m.Name == motion {
			return bakes.Frames / m.Length
		}
	}
	return 1
}
