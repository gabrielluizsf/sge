package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/vector"
)

// How near a cow the player must be to milk her; how long, in seconds, they
// take to go to the stool's place by her side, and how much of a cell is
// kept between that place and the fence.
const (
	milkReach  = 1.8 * tile
	milkGlide  = 0.35
	milkMargin = 0.2 * tile
)

// milkTotal is how long milking takes: sitting down, milking and getting up.
const milkTotal = 2*MilkDownSeconds + MilkSeconds

// milking is the player milking a cow: how long they have been at it, from
// where to where they go to sit, and the way they face.
type milking struct {
	c        *creature
	time     float64
	from, to vector.Vector3
	// cowFrom and cowTo are where the cow steps to keep the place in her pen.
	cowFrom, cowTo vector.Vector3
	dir            string
}

// cowBy is the cow nearest o, out of doors and within reach, unless the
// player is milking already.
func (f *farm) cowBy(o *game.Object) *creature {
	if f.herd == nil || f.room != nil || f.milking != nil || f.hands.swims {
		return nil
	}
	var near *creature
	best := milkReach
	for _, c := range f.herd.crits {
		if d := c.pos.Sub(feet(o)).Length(); c.kind == "cow" && c.pen.town == nil && d < best {
			near, best = c, d
		}
	}
	return near
}

// milkedToday reports whether c has given her bottle of the day.
func (f *farm) milkedToday(c *creature) bool { return c.milked == f.clock.Day+1 }

// cowNear is the interaction of a cow within reach: milk her, or be told
// she has been today.
func (f *farm) cowNear(o *game.Object) (words string, do func(), ok bool) {
	c := f.cowBy(o)
	switch {
	case c == nil:
		return "", nil, false
	case f.milkedToday(c):
		return i18n.T("Already milked today"), func() { f.sound("deny", nil) }, true
	}
	return i18n.T("E: milk the cow"), func() { f.milkNow(c) }, true
}

// milkNow starts the player milking c, shown at once: what E does.
func (f *farm) milkNow(c *creature) {
	if !f.startMilking(c) {
		return
	}
	clip, seconds := useClip(actMilk)
	f.hands.acting = seconds
	f.player.Animator.Play(f.milking.dir + clip)
}

// startMilking has the player sit by c, who stands still, turned so they
// can reach under her, and reports whether they did: a cow gives one
// bottle a day.
func (f *farm) startMilking(c *creature) bool {
	if f.milking != nil || f.milkedToday(c) {
		return false
	}
	cowFrom := c.pos
	dir, to := f.herd.milkSeat(c)
	c.walking, c.grazes, c.held = false, false, true
	c.heading = map[string]int{"right": 0, "left": 4}[dir]
	f.milking = &milking{c: c, from: feet(f.player), to: to, cowFrom: cowFrom, cowTo: c.pos, dir: dir}
	c.pos = cowFrom
	f.hands.facing = f.onGround(vector.Vector3{X: map[string]float64{"right": 1, "left": -1}[dir]})
	f.sound("ui_tick", nil)
	return true
}

// onGround is the way on the ground that the screen's way s shows.
func (f *farm) onGround(s vector.Vector3) vector.Vector3 {
	if g := f.player.Game(); g != nil && g.Projection != nil {
		return g.Projection.Unproject(s)
	}
	return s
}

// milkSeat is the way to face to milk c, and where to sit doing it, on the
// side that takes her less far from where she is to keep the place in her
// pen; she is moved that little.
func (h *herd) milkSeat(c *creature) (string, vector.Vector3) {
	back, near := h.art.Milker()
	dir, seat, shift := "", vector.Vector3{}, vector.Vector3{}
	for _, side := range []struct {
		dir  string
		sign float64
	}{{"right", -1}, {"left", 1}} {
		at := c.pos.Add(h.f.onGround(vector.Vector3{X: side.sign * back, Y: near}))
		moved := keepIn(c.pen.r.inside(), at).Sub(at)
		if dir == "" || moved.Length() < shift.Length() {
			dir, seat, shift = side.dir, at, moved
		}
	}
	c.pos = c.pos.Add(shift)
	c.to = c.pos
	return dir, seat.Add(shift)
}

// keepIn is at, moved the least to be inside the field r.
func keepIn(r rect, at vector.Vector3) vector.Vector3 {
	at.X = min(max(at.X, float64(r.col)*tile+milkMargin), float64(r.col+r.cols)*tile-milkMargin)
	at.Y = min(max(at.Y, float64(r.row)*tile+milkMargin), float64(r.row+r.rows)*tile-milkMargin)
	return at
}

// milkTick goes on with the milking: the player goes to their place, sits,
// milks and gets up, with the bottle.
func (f *farm) milkTick(dt float64) {
	m := f.milking
	if m == nil {
		return
	}
	m.time += dt
	u := min(m.time/milkGlide, 1)
	at := m.from.Add(m.to.Sub(m.from).Scale(u))
	f.player.Transform.Position.X, f.player.Transform.Position.Y = at.X, at.Y
	m.c.pos = m.cowFrom.Add(m.cowTo.Sub(m.cowFrom).Scale(u))
	m.c.to = m.c.pos
	switch {
	case m.time >= milkTotal:
		f.endMilking()
	case m.time >= MilkDownSeconds+MilkSeconds:
		f.player.Animator.Play(m.dir + "_milkup")
	case m.time >= MilkDownSeconds:
		f.player.Animator.Play(m.dir + "_milk")
	}
}

// endMilking is the bottle full: the player has it, and the cow goes on.
func (f *farm) endMilking() {
	m := f.milking
	m.c.milked = f.clock.Day + 1
	f.bag.Add(milk, 1)
	at := feet(f.player)
	at.Y -= tile
	f.float(at, "+1 "+i18n.T("Milk bottle"), 'y')
	f.sound("harvest", &at)
	f.earnBy(f.kit, xpMilk)
	f.releaseCow()
	f.hands.acting = 0
}

// stopMilking is the milking broken off, with no milk: a blow, or leaving.
func (f *farm) stopMilking() {
	if f.milking != nil {
		f.releaseCow()
	}
}

// releaseCow lets the cow milked go on with her day.
func (f *farm) releaseCow() {
	c := f.milking.c
	c.held, c.wait = false, 2
	f.milking = nil
}
