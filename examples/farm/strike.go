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
	clip := "_act"
	if p.kit.slot == slotAxe && p.kit.axe != noAxe {
		clip = "_chop"
	}
	p.acting = actFor
	o.Animator.Play(p.pose(o) + clip)
	at := feet(o)
	p.f.sound("chop", &at)
	if pk := p.f.wolves; pk != nil {
		if w := pk.ahead(at, p.facing); w != nil {
			if o == p.f.player {
				p.f.killedBy(pk, w, p.kit.blow(), at)
			} else {
				pk.strike(w, p.kit.blow(), at)
			}
			return
		}
	}
	if o == p.f.player {
		if v := p.f.ahead(at, p.facing); v != nil {
			p.f.wound(v, p.kit.blow(), at)
		}
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
