package farm

import (
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// person is who a person's object is: its name, how the art draws them,
// their kit, and what controls them.
type person struct {
	name string
	body Body
	kit  *kit
	ctl  Controller
}

// newFarmer is the player's farmer, whom the camera follows.
func newFarmer(f *farm, ctl Controller) (*game.Object, *farmer) {
	o, p := newPerson(f, person{name: "farmer", body: f.art.Person(f.self.who(f.clock.Day)), kit: f.kit, ctl: ctl})
	// The camera follows the farmer, easing after them.
	o.AddScript(game.Follow{Lag: 0.12})
	return o, p
}

// newPerson is a person's object, standing still and facing down, with
// their animations: standing, walking, pouring, chopping and swimming, in
// each way their body faces.
func newPerson(f *farm, who person) (*game.Object, *farmer) {
	o := game.NewObject(who.name)
	o.Z = zThings
	o.Collider = feetBox
	o.Animator = animation.NewAnimator()
	n := who.body.Facings()
	for _, dir := range Facings(n) {
		addMoves(o.Animator, who.body, dir)
	}
	p := &farmer{f: f, kit: who.kit, ctl: who.ctl, facing: vector.Vector3{Y: 1}, facings: n}
	o.AddScript(p)
	o.Animator.Play("down")
	return o, p
}

// addMoves gives an the clips of body facing dir, each under the name the
// farmer plays it by: dir standing, and dir with "_walk", "_act", "_chop"
// or "_swim" after it.
func addMoves(an *animation.Animator, body Body, dir string) {
	clip := func(anim string, loop bool) *animation.Animation {
		c := body.Clip(anim, dir)
		return &animation.Animation{Frames: c.Frames, FPS: c.FPS, Loop: loop}
	}
	an.Add(dir, clip(AnimStand, true))
	an.Add(dir+"_walk", clip(AnimWalk, true))
	an.Add(dir+"_act", clip(AnimAct, false))
	an.Add(dir+"_chop", clip(AnimChop, false))
	an.Add(dir+"_swim", clip(AnimSwim, true))
}
