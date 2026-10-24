package main

import (
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
)

// newFarmer is the farmer object, standing still and facing down, with
// their animations: standing, walking and using, in each direction.
func newFarmer(f *farm, ctl controller) (*game.Object, *farmer) {
	o := game.NewObject("farmer")
	o.Z = zFarmer
	o.Collider = feetBox
	o.Animator = animation.NewAnimator()
	frame := func(name string) *sprite.Sprite { return sprite.New(look("farmer_" + name)) }
	for dir := range directions {
		o.Animator.Add(dir, &animation.Animation{Frames: []*sprite.Sprite{frame(dir)}, FPS: 1, Loop: true})
		o.Animator.Add(dir+"_walk", &animation.Animation{
			Frames: []*sprite.Sprite{frame(dir + "1"), frame(dir), frame(dir + "2"), frame(dir)},
			FPS:    8,
			Loop:   true,
		})
		o.Animator.Add(dir+"_act", &animation.Animation{Frames: []*sprite.Sprite{frame(dir + "Act")}, FPS: 1})
	}
	p := &farmer{f: f, ctl: ctl, facing: "down"}
	o.AddScript(p)
	// The camera follows the farmer, easing after them.
	o.AddScript(game.Follow{Lag: 0.12})
	o.Animator.Play("down")
	return o, p
}
