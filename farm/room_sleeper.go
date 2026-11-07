package farm

import (
	"github.com/gabrielluizsf/sge/animation"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// asleep reports whether v, at home, is in bed: it is their night, and
// they are past babyhood, as the babies are kept out of sight.
func (f *farm) asleep(v *villager) bool {
	return v.now(f) == doSleep && v.stage() >= life.Child
}

// lies reports whether v is to be shown in a bed of r: they are asleep,
// the room has one free for them, and the art has them lying in it. It
// gives the berth.
func (f *farm) lies(v *villager, r *room) (berth, bool) {
	if !f.asleep(v) {
		return berth{}, false
	}
	at, ok := r.bedFor(v.partner != nil, f.bedOrder(v))
	if !ok {
		return berth{}, false
	}
	if v.body == nil {
		f.age(v, f.clock.Day)
		v.body = f.art.Person(v.who)
	}
	if !sleepsAs(v.body, r.beds[at.bed].anim(at.side)) {
		if at, ok = r.freeSide(false, [2]int{}); !ok || !sleeps(v.body) {
			return berth{}, false
		}
	}
	return at, true
}

// bedOrder is 0 for the man of a couple and 1 for the woman, who lie
// at the left and the right of their double bed.
func (f *farm) bedOrder(v *villager) int {
	if v.woman {
		return 1
	}
	return 0
}

// sleeper is v's object asleep at berth at of r: lying in the bed's middle
// when it is single; on their side of it when it is double, as the look's
// pictures are of that side, or a little aside the middle when it is not.
func (f *farm) sleeper(v *villager, r *room, at berth) *game.Object {
	b := &r.beds[at.bed]
	o := game.NewObject("villager")
	o.Z = zThings
	o.Animator = animation.NewAnimator()
	name := b.anim(at.side)
	c := v.body.Clip(name, "down")
	o.Animator.Add(name, &animation.Animation{Frames: c.Frames, FPS: c.FPS, Loop: true})
	o.Animator.Play(name)
	o.Transform.Position = b.center(f)
	if b.double && !sleepBaked(f.art) {
		o.Transform.Position = o.Transform.Position.Add(vector.Vector3{Y: float64(2*at.side-1) * -doubleAside * tile})
	}
	if sleepBaked(f.art) && !b.double {
		r.things.Set(b.c0, b.r0, tilemap.Empty)
	}
	b.used, b.taken[at.side] = true, true
	f.residents = append(f.residents, resident{v: v, obj: o})
	return o
}

// doubleAside is how far, in cells, a sleeper in a double bed lies off its
// middle, across it, when they are drawn alone.
const doubleAside = 7.0 / 16
