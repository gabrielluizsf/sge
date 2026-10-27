package main

import (
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// newFarmer is the farmer object, standing still and facing down, with
// their animations: standing, walking and using, in each direction.
func newFarmer(f *farm, ctl controller) (*game.Object, *farmer) {
	o := game.NewObject("farmer")
	o.Z = zThings
	o.Collider = feetBox
	o.Animator = animation.NewAnimator()
	// Each picture stands on the farmer's feet, its bottom on their place.
	frame := func(name string) *sprite.Sprite {
		s := sprite.New(picture(spriteSVG("farmer_"+name), art.FarmerW, art.FarmerH))
		s.PivotY = 1
		return s
	}
	for dir := range directions {
		o.Animator.Add(dir, &animation.Animation{Frames: []*sprite.Sprite{frame(dir)}, FPS: 1, Loop: true})
		o.Animator.Add(dir+"_walk", &animation.Animation{
			Frames: []*sprite.Sprite{frame(dir + "1"), frame(dir), frame(dir + "2"), frame(dir)},
			FPS:    8,
			Loop:   true,
		})
		o.Animator.Add(dir+"_act", &animation.Animation{Frames: []*sprite.Sprite{frame(dir + "Act")}, FPS: 1})
	}
	p := &farmer{f: f, ctl: ctl, facing: vector.Vector3{Y: 1}}
	o.AddScript(p)
	// The camera follows the farmer, easing after them.
	o.AddScript(game.Follow{Lag: 0.12})
	o.Animator.Play("down")
	return o, p
}
