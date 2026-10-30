package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// farmer is a person, the player or a villager, who walks as their
// controller says and uses what they hold
// on the cell they work: the one the mouse picks within reach, or the one
// they face. In water they swim, slower, and have nothing to work.
type farmer struct {
	f      *farm
	kit    *kit
	ctl    Controller
	facing vector.Vector3 // the way they face, in the world
	aim    *[2]int        // the cell the mouse picks, if any
	acting float64        // seconds left of the use being shown
	swims  bool           // whether they are in the water
	heard  heard          // what of them was last heard
	// facings is how many ways their body faces on screen, 4 or 8.
	facings int
}

const (
	walkSpeed = 60 * scale
	actFor    = 0.3
	// reach is how far from the farmer's feet a cell's middle may be for
	// the mouse to pick it.
	reach = 1.6 * tile
)

// directions are the ways a controller may turn the farmer to face, as
// steps across the world's cells.
var directions = map[string][2]int{"down": {0, 1}, "up": {0, -1}, "left": {-1, 0}, "right": {1, 0}}

func (p *farmer) Update(o *game.Object, dt float64) {
	if p.acting > 0 {
		p.acting -= dt
		return
	}
	in := p.ctl.control(p.f, o, dt)
	p.aim = in.aim
	if d, ok := directions[in.face]; ok {
		p.facing = vector.Vector3{X: float64(d[0]), Y: float64(d[1])}
	}
	moves := p.move(o, in, dt)
	p.swims = p.f.inWater(feet(o))
	switch {
	case p.swims:
		o.Animator.Play(p.pose(o) + "_swim")
	case moves:
		o.Animator.Play(p.pose(o) + "_walk")
	default:
		o.Animator.Play(p.pose(o))
	}
	p.footsteps(o)
	if in.use && !p.swims {
		col, row := p.target(o)
		if p.f.use(p.kit, col, row) {
			p.acting = actFor
			if p.kit.slot == slotAxe {
				o.Animator.Play(p.pose(o) + "_chop")
			} else {
				o.Animator.Play(p.pose(o) + "_act")
			}
		}
	}
}

// move moves the farmer as in says, for dt seconds, at the speed they go
// where they are, and reports whether they were told to.
func (p *farmer) move(o *game.Object, in intent, dt float64) bool {
	l := in.dir.Length()
	if l == 0 {
		return false
	}
	p.facing = in.dir.Scale(1 / l)
	step := walkSpeed * dt
	if p.swims {
		step = swimSpeed * dt
	}
	if in.upTo > 0 {
		step = min(step, in.upTo)
	}
	d := in.dir.Scale(step / l)
	tilemap.Move(o, d.X, d.Y, p.f.solids()...)
	return true
}

// pose is the way the farmer's picture faces: the screen's way they face,
// as the world's turns on it.
func (p *farmer) pose(o *game.Object) string {
	s := p.facing
	if pr := o.Game().Projection; pr != nil {
		s = pr.Project(s)
	}
	return facingOf(s.X, s.Y, p.facings)
}

// target is the cell the farmer works: the mouse's, or the one in front.
func (p *farmer) target(o *game.Object) (col, row int) {
	if p.aim != nil {
		return p.aim[0], p.aim[1]
	}
	at := feet(o).Add(p.facing.Scale(0.75 * tile))
	return p.f.cellOf(at)
}
