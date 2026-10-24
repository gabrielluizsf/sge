package main

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
)

// farmer walks the farm as its controller says, and uses what they hold
// on the cell they face.
type farmer struct {
	f      *farm
	ctl    controller
	facing string
	acting float64 // seconds left of the use being shown
}

const (
	walkSpeed = 60 * scale
	actFor    = 0.3
)

// directions are where the farmer can face, and the cell beside them that
// way.
var directions = map[string][2]int{"down": {0, 1}, "up": {0, -1}, "left": {-1, 0}, "right": {1, 0}}

func (p *farmer) Update(o *game.Object, dt float64) {
	if p.acting > 0 {
		p.acting -= dt
		return
	}
	in := p.ctl.control(p.f, o)
	if in.face != "" {
		p.facing = in.face
	}
	if l := in.dir.Length(); l > 0 {
		if math.Abs(in.dir.X) > math.Abs(in.dir.Y) {
			p.facing = map[bool]string{true: "right", false: "left"}[in.dir.X > 0]
		} else {
			p.facing = map[bool]string{true: "down", false: "up"}[in.dir.Y > 0]
		}
		step := walkSpeed * dt
		if in.upTo > 0 {
			step = min(step, in.upTo)
		}
		d := in.dir.Scale(step / l)
		tilemap.Move(o, d.X, d.Y, p.f.water, p.f.things, p.f.wild)
		o.Animator.Play(p.facing + "_walk")
	} else {
		o.Animator.Play(p.facing)
	}
	if in.use {
		col, row := p.target(o)
		if p.f.use(col, row) {
			p.acting = actFor
			o.Animator.Play(p.facing + "_act")
		}
	}
}

// target is the cell the farmer faces.
func (p *farmer) target(o *game.Object) (col, row int) {
	at := feet(o)
	col, row, _ = p.f.ground.CellAt(at.X, at.Y)
	d := directions[p.facing]
	return col + d[0], row + d[1]
}
