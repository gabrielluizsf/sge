package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// How fast a horse goes under a rider, in cells a second, walking and with
// Shift held: more than a person on foot. How far to the side of it one
// stands to get on, in cells; how far behind them it is drawn; how near a
// horse one must be.
const (
	rideWalk   = 5.2
	rideRun    = 7.6
	mountAside = 0.75
	heldBehind = 0.01
)

// riding is the player on a horse: mid getting on or down for as long as time
// has not reached MountSeconds, from and to the way their feet go meanwhile.
type riding struct {
	c        *creature
	time     float64
	changing bool
	down     bool
	moved    bool // the rider went since the horse last looked
	from, to vector.Vector3
}

// riding reports whether o is the player, seated on a horse and free to go.
func (f *farm) riding(o *game.Object) bool {
	return f.ride != nil && !f.ride.changing && o == f.player
}

// horseBy is the player's horse nearest o, within reach, out of doors.
func (f *farm) horseBy(o *game.Object) *creature {
	if f.herd == nil || f.room != nil || f.milking != nil || f.ride != nil || f.hands.swims {
		return nil
	}
	var near *creature
	best := milkReach
	for _, c := range f.herd.crits {
		if d := c.pos.Sub(feet(o)).Length(); c.kind == "horse" && c.pen.town == nil && d < best {
			near, best = c, d
		}
	}
	return near
}

// horseNear is the interaction of a horse within reach: get on it.
func (f *farm) horseNear(o *game.Object) (words string, do func(), ok bool) {
	c := f.horseBy(o)
	if c == nil {
		return "", nil, false
	}
	return i18n.T("E: ride the horse"), func() { f.getOn(c) }, true
}

// besideHorse is the spot on the horse's left, as it faces, where one stands
// to get on or comes down: the way it goes being the screen's, its heading.
func (f *farm) besideHorse(c *creature) vector.Vector3 {
	at := c.pos.Add(f.onGround(leftOf(c.heading)).Scale(mountAside * tile))
	if col, row := f.cellOf(at); f.open(col, row) {
		return at
	}
	return c.pos.Add(f.onGround(leftOf(c.heading)).Scale(-mountAside * tile))
}

// leftOf is a unit way on the screen to the left of heading, an eighth of a
// turn each, clockwise from the right.
func leftOf(heading int) vector.Vector3 {
	return headingWay(heading + 6)
}

// headingWay is the way on the screen an eighth-of-a-turn heading points.
func headingWay(heading int) vector.Vector3 {
	ways := [8][2]float64{{1, 0}, {0.894, 0.447}, {0, 1}, {-0.894, 0.447}, {-1, 0}, {-0.894, -0.447}, {0, -1}, {0.894, -0.447}}
	w := ways[(heading%8+8)%8]
	return vector.Vector3{X: w[0], Y: w[1]}
}

// getOn starts the player getting on horse c, which stands for it, and
// reports whether they did.
func (f *farm) getOn(c *creature) bool {
	if f.ride != nil || f.seat != nil {
		return false
	}
	f.face8(c.heading)
	c.walking, c.grazes, c.held, c.homing = false, false, true, false
	f.ride = &riding{c: c, changing: true, from: f.besideHorse(c), to: c.pos}
	f.setFeet(f.ride.from)
	f.hands.acting = MountSeconds
	f.player.Animator.Play(f.hands.pose(f.player) + "_mount")
	f.sound("ui_tick", nil)
	return true
}

// face8 turns the player the way heading looks on the screen.
func (f *farm) face8(heading int) {
	f.hands.facing = unitOf(f.onGround(headingWay(heading)), f.hands.facing)
}

// setFeet stands the player at, where it is, as high as they are.
func (f *farm) setFeet(at vector.Vector3) {
	f.player.Transform.Position.X, f.player.Transform.Position.Y = at.X, at.Y
}

// getDown starts the player down from the horse.
func (f *farm) getDown() {
	r := f.ride
	if r == nil || r.changing {
		return
	}
	r.c.walking = false
	r.changing, r.down, r.time = true, true, 0
	r.from, r.to = feet(f.player), f.besideHorse(r.c)
	f.hands.acting = MountSeconds
	f.player.Animator.Play(f.hands.pose(f.player) + "_dismount")
	f.sound("ui_tick", nil)
}

// rideTick goes on with getting on or down, or keeps the horse under the rider.
func (f *farm) rideTick(dt float64) {
	r := f.ride
	if r == nil {
		return
	}
	if !r.changing {
		r.c.pos, r.c.to = feet(f.player), feet(f.player)
		r.c.walking, r.moved = r.moved, false
		return
	}
	r.time += dt
	u := min(r.time/MountSeconds, 1)
	f.setFeet(r.from.Add(r.to.Sub(r.from).Scale(u)))
	if u < 1 {
		return
	}
	f.hands.acting = 0
	r.changing = false
	if r.down {
		f.freeHorse(r.c)
		f.ride = nil
	}
}

// throwOff is the rider gone at once, dead: the horse goes home by itself.
func (f *farm) throwOff() {
	if f.ride != nil {
		f.freeHorse(f.ride.c)
		f.ride = nil
	}
}

// freeHorse sends c, its rider down, back to its pen on its own.
func (f *farm) freeHorse(c *creature) {
	c.held, c.walking, c.mounted = false, false, false
	c.homing, c.route = true, f.herd.wayHome(c)
	if c.route == nil { // too far to find the way: straight for its pen
		c.to = f.herd.spot(c.pen)
	}
}

// ride moves the player on their horse for dt seconds: over land, at the
// horse's speed, which stops at deep water; and it walks while it goes.
func (p *farmer) ride(o *game.Object, in intent, dt float64) bool {
	speed := float64(rideWalk)
	if in.run {
		speed = rideRun
	}
	r := p.f.ride
	r.moved, r.c.mounted, r.c.rideSpeed = true, true, speed
	r.c.walkSpeed, _ = p.f.herd.art.Pace("horse")
	r.c.heading = p.f.herd.heading(p.facing)
	d := in.dir.Scale(speed * tile * dt / in.dir.Length())
	before := feet(o)
	tilemap.Move(o, d.X, d.Y, p.f.solids()...)
	if p.f.inWater(feet(o)) {
		o.Transform.Position.X, o.Transform.Position.Y = before.X, before.Y
	}
	return true
}

// pace is how many times faster than its own walk a creature's legs go,
// under a rider, for its feet not to slide.
func (c *creature) pace() float64 {
	if !c.mounted || c.walkSpeed <= 0 {
		return 1
	}
	return c.rideSpeed / c.walkSpeed
}

// wayHome is the way from where c is to its pen's gate and in, or nil when
// there is none.
func (h *herd) wayHome(c *creature) *pathfind.Route {
	f := h.f
	p := c.pen
	col, row := f.cellOf(c.pos)
	inside := pathfind.Cell{Col: p.r.col + p.r.cols/2, Row: p.r.row + p.r.rows - 2}
	way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, inside)
	if len(way) == 0 {
		return nil
	}
	return &pathfind.Route{Cells: way}
}

// goHome walks c along its way to its pen, or straight there with none
// found, and there it lives as before.
func (h *herd) goHome(c *creature, dt float64) {
	way, done := c.to.Sub(c.pos), false
	if c.route != nil {
		way, done = c.route.Steer(c.pos, h.f.center)
	} else if way.Length() < 1 {
		done = true
	}
	if done {
		c.homing, c.route, c.walking = false, nil, false
		c.wait = 2
		return
	}
	speed, _ := h.art.Pace(c.kind)
	c.walking = true
	step := min(speed*tile*dt, way.Length())
	c.pos = c.pos.Add(way.Scale(step / way.Length()))
	if c.route != nil {
		c.to = c.pos
	}
	c.heading = h.heading(way)
}

// backInPen puts c at once in its pen: no one is looking.
func (h *herd) backInPen(c *creature) {
	c.pos = h.spot(c.pen)
	c.to = c.pos
	c.homing, c.route, c.walking = false, nil, false
	c.wait = 2
}
