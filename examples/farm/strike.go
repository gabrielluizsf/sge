package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/vector"
)

// swing is the farmer swinging what they hold at whoever is in front of
// them: the stroke of the axe if it is in hand, else of the bare arm. A
// villager's blow lands on the player, the player's on the villager in
// front of them.
func (p *farmer) swing(o *game.Object) {
	if o == p.f.player && p.f.wolves != nil {
		if w := p.f.wolves.nearest(feet(o), wolfBlow); w != nil {
			p.facing = unitOf(feet(w.obj).Sub(feet(o)), p.facing) // a blow turns them to it
		}
	}
	clip := "_act"
	if p.kit.slot == slotAxe && p.kit.axe != noAxe {
		clip = "_chop"
	}
	p.acting = actFor
	o.Animator.Play(p.pose(o) + clip)
	at := feet(o)
	p.f.sound("chop", &at)
	if pk := p.f.wolves; pk != nil {
		if o == p.f.player {
			if w := pk.nearest(at, wolfBlow); w != nil {
				p.f.killedBy(pk, w, p.kit.blow(), at)
				return
			}
		} else if w := pk.ahead(at, p.facing); w != nil {
			pk.strike(w, p.kit.blow(), at)
			return
		}
	}
	if o == p.f.player {
		if v := p.f.ahead(at, p.facing); v != nil {
			p.f.wound(v, p.kit.blow(), at)
			return
		}
		p.f.strikeAnimal(at, p.facing, p.kit.blow())
		return
	}
	if m, ok := p.ctl.(*mind); ok && m.v.angry <= 0 {
		return // swinging at a wolf, not at the player
	}
	if p.f.player.Transform.Position.Sub(at).Length() <= reachBlow {
		p.f.hurtPlayer(npcHarm, at)
	}
}

// ahead is the villager nearest in front of at, facing the way facing
// points, within reach of a blow, or nil.
func (f *farm) ahead(at, facing vector.Vector3) *villager {
	var best *villager
	far := reachBlow
	for _, v := range f.people {
		if v.obj == nil || v.obj.Hidden || v.inside || v.ours || v.stage() < life.Teen {
			continue
		}
		d := feet(v.obj).Sub(at)
		l := d.Length()
		if l > far || l > 0 && d.Dot(facing)/l < 0.3 {
			continue
		}
		best, far = v, l
	}
	return best
}

// unitOf is d made a single unit long, or else when it has no length.
func unitOf(d, otherwise vector.Vector3) vector.Vector3 {
	if l := d.Length(); l > 0 {
		return d.Scale(1 / l)
	}
	return otherwise
}

// fightsWolf reports whether a use of what p holds is a blow at a wolf: the
// player with a tool that strikes in hand, a wolf within reach.
func (p *farmer) fightsWolf(o *game.Object) bool {
	if o != p.f.player || p.f.wolves == nil || p.kit.slot != slotAxe && p.kit.slot != slotPick {
		return false
	}
	return p.f.wolves.nearest(feet(o), wolfBlow) != nil
}
