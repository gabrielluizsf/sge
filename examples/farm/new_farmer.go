package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// person is who a person's object is: its name, the SVG of each of its
// frames by name, their kit, and what controls them.
type person struct {
	name  string
	frame func(name string) string
	kit   *kit
	ctl   controller
}

// newFarmer is the player's farmer, whom the camera follows.
func newFarmer(f *farm, ctl controller) (*game.Object, *farmer) {
	o, p := newPerson(f, person{name: "farmer", frame: func(n string) string { return spriteSVG("farmer_" + n) }, kit: f.kit, ctl: ctl})
	// The camera follows the farmer, easing after them.
	o.AddScript(game.Follow{Lag: 0.12})
	return o, p
}

// newPerson is a person's object, standing still and facing down, with
// their animations: standing, walking, pouring, chopping and swimming, in
// each direction.
func newPerson(f *farm, who person) (*game.Object, *farmer) {
	o := game.NewObject(who.name)
	o.Z = zThings
	o.Collider = feetBox
	o.Animator = animation.NewAnimator()
	// Each picture stands on the person's feet, its bottom on their place,
	// drawn at half the size of its art when it is first shown.
	frame := func(name string) *sprite.Sprite {
		return &sprite.Sprite{PivotX: 0.5, PivotY: 1, Make: func() *canvas.Canvas {
			return picture(who.frame(name), art.FarmerW/2, art.FarmerH/2)
		}}
	}
	for dir := range directions {
		o.Animator.Add(dir, &animation.Animation{Frames: []*sprite.Sprite{frame(dir)}, FPS: 1, Loop: true})
		o.Animator.Add(dir+"_walk", &animation.Animation{
			Frames: []*sprite.Sprite{frame(dir + "1"), frame(dir), frame(dir + "2"), frame(dir)},
			FPS:    8,
			Loop:   true,
		})
		o.Animator.Add(dir+"_act", &animation.Animation{Frames: []*sprite.Sprite{frame(dir + "Act")}, FPS: 1})
		o.Animator.Add(dir+"_chop", &animation.Animation{Frames: []*sprite.Sprite{frame(dir + "Chop"), frame(dir + "Chop1")}, FPS: 2 / actFor})
		o.Animator.Add(dir+"_swim", &animation.Animation{Frames: []*sprite.Sprite{frame(dir + "Swim"), frame(dir + "Swim1")}, FPS: 3, Loop: true})
	}
	p := &farmer{f: f, kit: who.kit, ctl: who.ctl, facing: vector.Vector3{Y: 1}}
	o.AddScript(p)
	o.Animator.Play("down")
	return o, p
}
