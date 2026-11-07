package farm

import "github.com/gabrielluizsf/sge/game"

// cross shows the farmer going into deep water, where p.swims has just
// become true: a dive, or out of it, where it has become false: standing
// up. Each is shown once, the farmer still while it is.
func (p *farmer) cross(o *game.Object) {
	clip, seconds := "_rise", RiseSeconds
	if p.swims {
		clip, seconds = "_dive", DiveSeconds
	}
	o.Animator.Play(p.pose(o) + clip)
	p.acting = seconds
	p.footsteps(o)
}
