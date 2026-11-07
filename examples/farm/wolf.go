package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// What a wolf is and does: its health, how far it bites, how hard and how
// often; how far off it comes from, how far it hunts, how fast it runs
// against its walk, how long it lies dead; how often the way to its prey
// is found anew.
const (
	wolfHealth = 45.0
	wolfReach  = 1.1 * tile
	wolfBite   = 8.0
	wolfEvery  = 1.5
	wolfFar    = 15 * tile
	wolfHunts  = 26 * tile
	wolfRun    = 2.2
	wolfDead   = 8.0
	wolfThink  = 0.5
	biteFor    = 0.5
)

// Wolves come from dusk until dawn: at most packMost at once, one every
// packEvery seconds or so.
const (
	duskAt    = 18.0
	dawnAt    = 6.0
	packMost  = 3
	packEvery = 25.0
)

// wolf is a wolf about: where it is, hurt how much, the way it faces on
// the screen, the time in its motion, its next bite and its way.
type wolf struct {
	obj     *game.Object
	wounds  float64
	heading int
	time    float64
	biting  float64
	next    float64
	think   float64
	route   *pathfind.Route
	dead    float64 // seconds it has lain dead, once it is
	fallen  bool
}

// wolfPack is the wolves about, seen to each frame from dusk until dawn. Those
// about at dawn are gone with the night, and none is saved: the next night
// brings others.
type wolfPack struct {
	f      *farm
	art    AnimalArt
	luck   *random.Source
	wolves []*wolf
	due    float64
	// on says wolves come: only to one playing, not to the demo's bot,
	// unless it is told so.
	on bool
}

// newWolves is the farm's pack, or nil when its art has no animals.
func newWolves(f *farm, on bool) *wolfPack {
	art, ok := f.art.(AnimalArt)
	if !ok {
		return nil
	}
	return &wolfPack{f: f, art: art, luck: random.New(seed + 61), on: on, due: 3}
}

// wolfHour reports whether hour is one wolves come at.
func wolfHour(hour float64) bool { return hour >= duskAt || hour < dawnAt }

// Update lets a new wolf come, when it is time, and each wolf live.
func (p *wolfPack) Update(o *game.Object, dt float64) {
	f := p.f
	night := wolfHour(f.clock.Hour)
	if !night {
		p.clear()
	}
	if p.on && night && f.room == nil {
		if p.due -= dt; p.due <= 0 {
			p.due = packEvery * (0.6 + 0.8*p.luck.Float())
			p.come(o.Game())
		}
	}
	kept := p.wolves[:0]
	for _, w := range p.wolves {
		w.obj.Hidden = f.room != nil
		if f.room == nil {
			p.live(w, dt)
		}
		if w.fallen && w.dead > wolfDead {
			w.obj.Destroy()
			continue
		}
		kept = append(kept, w)
	}
	p.wolves = kept
}

// clear is the night over: the wolves go.
func (p *wolfPack) clear() {
	for _, w := range p.wolves {
		w.obj.Destroy()
	}
	p.wolves = p.wolves[:0]
}

// come brings a wolf from out of sight, if there is ground for it.
func (p *wolfPack) come(g *game.Game) {
	live := 0
	for _, w := range p.wolves {
		if !w.fallen {
			live++
		}
	}
	if live >= packMost {
		return
	}
	f := p.f
	at := feet(f.player)
	for range 12 {
		a := p.luck.Float() * 2 * math.Pi
		to := at.Add(vector.Vector3{X: wolfFar * math.Cos(a), Y: wolfFar * math.Sin(a)})
		col, row := f.cellOf(to)
		if f.open(col, row) {
			p.add(g, f.center(pathfind.Cell{Col: col, Row: row}))
			return
		}
	}
}

// add puts a wolf at pos.
func (p *wolfPack) add(g *game.Game, pos vector.Vector3) *wolf {
	o := game.NewObject("animal")
	o.Z = zThings
	o.Collider = feetBox
	o.Transform.Position = pos
	g.Add(o)
	w := &wolf{obj: o, next: 0.5}
	p.wolves = append(p.wolves, w)
	return w
}

// prey is what w goes for: the nearest, within its hunt, of the player out
// of doors and the villagers out in sight.
type prey struct {
	at vector.Vector3
	v  *villager // nil for the player
	ok bool
}

// preyOf is the prey nearest w.
func (p *wolfPack) preyOf(w *wolf) prey {
	f := p.f
	at := feet(w.obj)
	best := prey{}
	far := float64(wolfHunts)
	if d := feet(f.player).Sub(at).Length(); d < far && !f.hands.swims {
		best, far = prey{at: feet(f.player), ok: true}, d
	}
	for _, v := range f.people {
		if v.obj == nil || v.obj.Hidden || v.inside || v.ours || v.stage() < teenAndUp {
			continue
		}
		if d := feet(v.obj).Sub(at).Length(); d < far {
			best, far = prey{at: feet(v.obj), v: v, ok: true}, d
		}
	}
	return best
}

// live moves w: lying dead, biting what it has reached, or running to it.
func (p *wolfPack) live(w *wolf, dt float64) {
	w.time += dt
	if w.fallen {
		w.dead += dt
		p.dress(w)
		return
	}
	w.biting -= dt
	w.next -= dt
	target := p.preyOf(w)
	if target.ok {
		p.hunt(w, target, dt)
	}
	p.dress(w)
}

// hunt has w run to target, and bite it when it is near enough.
func (p *wolfPack) hunt(w *wolf, target prey, dt float64) {
	f := p.f
	at := feet(w.obj)
	d := target.at.Sub(at)
	if d.Length() <= wolfReach {
		w.heading = p.headingOf(d)
		w.route = nil
		if w.next <= 0 {
			w.next, w.biting, w.time = wolfEvery, biteFor, 0
			p.bite(target, at)
		}
		return
	}
	if w.biting > 0 {
		return
	}
	if w.think -= dt; w.route == nil || w.think <= 0 {
		w.think = wolfThink
		col, row := f.cellOf(at)
		tc, tr := f.cellOf(target.at)
		if way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, pathfind.Cell{Col: tc, Row: tr}); len(way) > 0 {
			w.route = &pathfind.Route{Cells: way}
		} else {
			w.route = nil
		}
	}
	way := d
	if w.route != nil {
		if s, done := w.route.Steer(at, f.center); !done {
			way = s
		}
	}
	speed, _ := p.art.Pace("wolf")
	step := min(speed*wolfRun*tile*dt, way.Length())
	if l := way.Length(); l > 0 {
		m := way.Scale(step / l)
		tilemap.Move(w.obj, m.X, m.Y, f.solids()...)
		w.heading = p.headingOf(way)
	}
}

// bite is the wolf at from biting target.
func (p *wolfPack) bite(target prey, from vector.Vector3) {
	if target.v != nil {
		p.f.bitten(target.v, wolfBite, from)
		return
	}
	p.f.hurtPlayer(wolfBite, from)
}

// headingOf is the eighth of a turn, clockwise from the screen's right,
// that one going along dir on the ground faces on the screen.
func (p *wolfPack) headingOf(dir vector.Vector3) int {
	s := dir
	if g := p.f.player.Game(); g != nil && g.Projection != nil {
		s = g.Projection.Project(dir)
	}
	return (int(math.Round(math.Atan2(s.Y, s.X)/(math.Pi/4)))%8 + 8) % 8
}

// dress gives w's object the picture of what it does.
func (p *wolfPack) dress(w *wolf) {
	_, fps := p.art.Pace("wolf")
	look := AnimalLook{Kind: "wolf", Heading: w.heading, Motion: "walk", Frame: int(w.time * fps * wolfRun)}
	switch {
	case w.fallen:
		look.Motion = "fall"
		look.Frame = min(int(w.time*fps), p.art.Frames("wolf", "fall")-1)
	case w.biting > 0:
		look.Motion = "bite"
		look.Frame = min(int((biteFor-w.biting)/biteFor*float64(p.art.Frames("wolf", "bite"))), p.art.Frames("wolf", "bite")-1)
	}
	if s := p.art.Animal(look); s != nil && s.Image != nil {
		w.obj.Sprite = s
	}
}

// ahead is the live wolf nearest in front of at, facing the way facing
// points, within reach of a blow, or nil.
func (p *wolfPack) ahead(at, facing vector.Vector3) *wolf {
	var best *wolf
	far := float64(reachBlow)
	for _, w := range p.wolves {
		if w.fallen {
			continue
		}
		d := feet(w.obj).Sub(at)
		l := d.Length()
		if l > far || l > 0 && d.Dot(facing)/l < 0.3 {
			continue
		}
		best, far = w, l
	}
	return best
}

// nearest is the live wolf nearest at within reach, or nil.
func (p *wolfPack) nearest(at vector.Vector3, reach float64) *wolf {
	var best *wolf
	for _, w := range p.wolves {
		if d := feet(w.obj).Sub(at).Length(); !w.fallen && d < reach {
			best, reach = w, d
		}
	}
	return best
}

// strike is w hit for harm by someone at from: it bleeds, and at the end of
// its health falls, and the player has the credit if they gave the blow.
func (p *wolfPack) strike(w *wolf, harm float64, from vector.Vector3) {
	f := p.f
	at := feet(w.obj)
	f.bleed(at, from, 14)
	if w.wounds += harm; w.wounds >= wolfHealth {
		w.fallen, w.time, w.dead = true, 0, 0
		w.obj.Collider = nil
		f.bleed(at, from, 30)
		f.sound("fell", &at)
	}
}
