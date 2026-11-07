package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
	"github.com/gabrielluizsf/sge/vitals"
)

// fadeFor is how many seconds a cut between shots fades out and in.
const fadeFor = 0.6

// shot is a scene of the film behind the main menu: how long it lasts, the
// hour it is set at, what is made ready for it, and what the farmer does.
type shot struct {
	seconds, hour float64
	ready         func(c *cinema, f *farm, o *game.Object)
	act           func(c *cinema, f *farm, o *game.Object, dt float64) intent
}

// cinema is the farmer of the main menu's background, a film of a farmer's
// life cut together: tending the field, felling trees, milking, trading in
// town, breaking rocks, and at night fighting off wolves with the villagers,
// a new day each time round.
type cinema struct {
	bot   bot
	at    int
	time  float64
	began bool
	// route is the way being walked, and next how long until the next
	// stroke; work the job being done.
	route *pathfind.Route
	next  float64
	work  *job
	// met is the villager being spoken to.
	met *villager
}

// Cinema is the main menu's farmer.
func Cinema() Controller { return &cinema{bot: bot{tends: 10}} }

// film is the shots in order.
var film = []shot{
	{seconds: 16, hour: 8, ready: readyField, act: actField},
	{seconds: 15, hour: 10.5, ready: readyChop, act: actStrokes},
	{seconds: milkTotal + 2, hour: 7, ready: readyMilk, act: actStill},
	{seconds: 14, hour: 14, ready: readyTrade, act: actTrade},
	{seconds: 13, hour: 16, ready: readyMine, act: actStrokes},
	{seconds: 18, hour: 19.5, ready: readyWolves, act: actWolves},
}

// Update keeps the film's time, which runs on while the farmer is busy.
func (c *cinema) Update(_ *game.Object, dt float64) {
	if c.began {
		c.time += dt
	}
}

func (c *cinema) control(f *farm, o *game.Object, dt float64) intent {
	f.vit, f.flash = vitals.Full(), 0
	f.talk, f.question = nil, nil
	s := film[c.at]
	if !c.began {
		c.began, c.time = true, 0
		c.route, c.work, c.met, c.next = nil, nil, nil, 0.5
		f.clock.Hour = s.hour
		s.ready(c, f, o)
		game.Follow{}.Start(o)
	}
	if c.time >= s.seconds {
		c.cut(f)
		return intent{}
	}
	return s.act(c, f, o, dt)
}

// cut goes on to the next shot; after the last, it is the next day.
func (c *cinema) cut(f *farm) {
	f.stopMilking()
	if f.wolves != nil {
		f.wolves.clear()
	}
	c.began = false
	if c.at = (c.at + 1) % len(film); c.at == 0 {
		f.passDays(1)
	}
}

// fade is how dark the screen is at this moment of the shot, 0 to 1: black
// at a cut, clear between.
func (c *cinema) fade() float64 {
	if !c.began {
		return 1
	}
	end := film[c.at].seconds - c.time
	return max(0, 1-min(c.time, end)/fadeFor)
}

// place stands o in the middle of cell at.
func place(f *farm, o *game.Object, at pathfind.Cell) {
	pos := f.center(at)
	o.Transform.Position.X, o.Transform.Position.Y = pos.X, pos.Y
}

// walk has o go along the route to c's job, reporting whether they are there.
func (c *cinema) walk(f *farm, o *game.Object) (intent, bool) {
	if c.route == nil {
		return intent{}, true
	}
	way, done := c.route.Steer(feet(o), f.center)
	if done {
		c.route = nil
		return intent{}, true
	}
	return intent{dir: way, upTo: way.Length()}, false
}

// goTo sets c walking from o to cell to, if there is a way.
func (c *cinema) goTo(f *farm, o *game.Object, to pathfind.Cell) {
	col, row := f.cellOf(feet(o))
	if way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, to); len(way) > 0 {
		c.route = &pathfind.Route{Cells: way}
	}
}

func readyField(c *cinema, f *farm, o *game.Object) {
	c.bot.task = nil
	place(f, o, pathfind.Cell{Col: stallCol, Row: stallRow + 1})
}

func actField(c *cinema, f *farm, o *game.Object, dt float64) intent {
	return c.bot.control(f, o, dt)
}

// quarryJobs are the places to strike the things is says, about the town
// next door.
func quarryJobs(f *farm, q quarry) []*job {
	t := f.towns.Of(1, 0)
	f.settle(t)
	q.area = area{center: pathfind.Cell(t.Center), cols: 15, rows: 9}
	q.where = func(col, row int) bool { return !f.onFarm(col, row) }
	return f.sitesOf(q, false)
}

func readyChop(c *cinema, f *farm, o *game.Object) {
	f.axe = max(f.axe, 0)
	f.slot = slotAxe
	jobs := quarryJobs(f, quarry{is: isTree, slot: slotAxe, uses: axes[f.axe].hits})
	if j := oneOf(jobs, uint32(f.clock.Day)); j != nil {
		c.work = j
		place(f, o, j.to)
	}
}

func readyMine(c *cinema, f *farm, o *game.Object) {
	if f.bag.Count(itemPick) == 0 {
		f.bag.Add(itemPick, 1)
	}
	f.slot = slotPick
	jobs := quarryJobs(f, quarry{is: isRock, slot: slotPick, uses: rockHits})
	if j := oneOf(jobs, uint32(f.clock.Day)); j != nil {
		c.work = j
		place(f, o, j.to)
	}
}

// actStrokes is striking the thing in front, a stroke at a time.
func actStrokes(c *cinema, f *farm, o *game.Object, dt float64) intent {
	if c.work == nil {
		return intent{}
	}
	if c.next -= dt; c.next > 0 {
		return intent{face: c.work.face}
	}
	c.next = 1.1
	return intent{face: c.work.face, use: true}
}

func readyMilk(c *cinema, f *farm, o *game.Object) {
	if f.herd == nil || len(f.herd.crits) == 0 {
		return
	}
	cow := f.herd.own("cow")
	if cow == nil {
		return
	}
	cow.milked = 0
	f.player.Transform.Position.X, f.player.Transform.Position.Y = cow.pos.X+tile, cow.pos.Y
	f.milkNow(cow)
}

func actStill(*cinema, *farm, *game.Object, float64) intent { return intent{} }

func readyTrade(c *cinema, f *farm, o *game.Object) {
	t := f.towns.Of(1, 0)
	f.settle(t)
	place(f, o, pathfind.Cell{Col: t.Stall.Col + 3, Row: t.Stall.Row + 3})
	c.goTo(f, o, pathfind.Cell{Col: t.Stall.Col, Row: t.Stall.Row + 1})
	f.bag.Add(wood, 4)
}

// actTrade is walking to the stall, selling there, and turning to whoever
// is nearest, who stops to talk.
func actTrade(c *cinema, f *farm, o *game.Object, dt float64) intent {
	in, there := c.walk(f, o)
	if !there {
		return in
	}
	if c.next -= dt; c.next > 0 {
		return intent{}
	}
	c.next = 3
	if worth(f.kit) > 0 {
		f.sell(f.kit, feet(o))
		return intent{face: "up"}
	}
	if v := f.nearest(feet(o), 8*tile); v != nil {
		c.met = v
		v.mind.hold, v.mind.facing = 3, faceTo(feet(v.obj), f.player)
		col, row := f.cellOf(feet(v.obj))
		c.goTo(f, o, pathfind.Cell{Col: col, Row: row + 1})
	}
	return intent{}
}

// nearest is the villager out nearest at, within reach, or nil.
func (f *farm) nearest(at vector.Vector3, reach float64) *villager {
	var best *villager
	for _, v := range f.people {
		if v.obj == nil || v.obj.Hidden || v.inside || v.ours {
			continue
		}
		if d := feet(v.obj).Sub(at).Length(); d < reach {
			best, reach = v, d
		}
	}
	return best
}

func readyWolves(c *cinema, f *farm, o *game.Object) {
	t := f.towns.Of(1, 0)
	f.settle(t)
	f.axe = max(f.axe, 1)
	f.slot = slotAxe
	place(f, o, pathfind.Cell{Col: t.Center.Col + 3, Row: t.Center.Row})
	if f.wolves != nil {
		c.loose(f, o)
	}
}

// loose lets two wolves in from either side of o.
func (c *cinema) loose(f *farm, o *game.Object) {
	for _, side := range []float64{-1, 1} {
		f.wolves.add(o.Game(), feet(o).Add(vector.Vector3{X: side * 7 * tile, Y: -side * 2 * tile}))
	}
}

// actWolves is going for the nearest wolf, and swinging once near; more
// come as those are felled.
func actWolves(c *cinema, f *farm, o *game.Object, dt float64) intent {
	if f.wolves == nil {
		return intent{}
	}
	w := f.wolves.nearest(feet(o), 20*tile)
	if w == nil {
		c.loose(f, o)
		return intent{}
	}
	d := feet(w.obj).Sub(feet(o))
	c.next -= dt
	if d.Length() <= reachBlow*0.9 {
		in := intent{dir: d, upTo: 0.01}
		if c.next <= 0 {
			c.next, in.strike = 0.7, true
		}
		return in
	}
	return intent{dir: d, run: true}
}
