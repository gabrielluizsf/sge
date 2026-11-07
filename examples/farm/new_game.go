package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/scene"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/warm"
)

// build is the farm laid out, drawn by a, with the farmer on the path in
// front of the stall, played by ctl; the towns around it, waking as the
// camera nears them; the clock, the night's shade over it all, and the
// HUD.
func build(a Art, ctl Controller) (*game.Game, *farm) {
	g := game.New()
	g.Background = canvas.RGB(43, 31, 46)
	// Seen from a corner: each cell of the world a diamond twice as wide
	// as it is tall.
	g.Projection = iso.New(tile, isoW, tile)
	// The farm and the land around it cover every pixel.
	g.Covered = true
	// Villagers far off the screen are not drawn yet: their pictures wait.
	g.LazyReach = 512
	tiles := newTiles(a)
	roomTiles = tiles
	where := siteFor(a)
	l := newLayers(g, tiles, where)
	towns := newTowns(a)
	wild := newWild(world{art: a, tiles: tiles, towns: towns}, [2]int{l.ground.Cols, l.ground.Rows})
	f := newFarm(l, world{art: a, tiles: tiles, wild: wild, towns: towns})
	w := game.NewObject("world")
	w.AddScript(f.clock)
	w.AddScript(f)
	f.warm = &warm.Queue{PerFrame: 1}
	w.AddScript(f.warm)
	g.Add(w)
	f.ui = newUI(f)
	menus := game.NewObject("menus")
	menus.Always = true
	menus.AddScript(f.ui)
	g.Add(menus)
	o, p := newFarmer(f, ctl)
	f.player, f.hands = o, p
	g.Reveals = f.reveals
	if bs := newBirds(f); bs != nil {
		w.AddScript(bs)
	}
	if hd := newHerd(f); hd != nil {
		f.herd = hd
		w.AddScript(hd)
	}
	if c, ok := ctl.(*cinema); ok {
		w.AddScript(c)
	}
	_, playing := ctl.(keyboard)
	if wv := newWolves(f, playing); wv != nil {
		f.wolves = wv
		w.AddScript(wv)
	}
	o.Transform.Position = standIn(f, where.spawn.Col, where.spawn.Row)
	g.Add(o)
	homeDoor = where.door
	g.Add(f.newStall(farmStall))
	g.Add(f.newHouse(0, town.Cell{Col: where.home.col, Row: where.home.row}, true))
	wild.Props = f.props
	wild.Warm = f.warm
	// The wild makes its chunks, and the towns wake, where the camera has
	// just been pointed, so they come after the farmer, whom it follows.
	g.Add(wild.Object)
	f.keeper = &town.Keeper{Layout: towns, Range: 34, Cell: tile, Wake: f.wake, Sleep: f.dropBodies, Ahead: 25, PerFrame: 1}
	k := game.NewObject("keeper")
	k.AddScript(f.keeper)
	g.Add(k)
	night := game.NewObject("night")
	night.Z, night.Screen = zNight, true
	night.AddScript(daytime.Shade{Clock: f.clock})
	g.Add(night)
	blood, stains := f.newBlood()
	g.Add(blood)
	g.Add(stains)
	h := game.NewObject("hud")
	h.Z, h.Screen = zHUD, true
	h.AddScript(&hud{f: f, farmer: o, p: p})
	g.Add(h)
	kept := map[*game.Object]bool{w: true, menus: true, o: true, night: true, h: true, blood: true, stains: true}
	f.stage = &scene.Stage{Game: g, Current: outside, Keep: func(o *game.Object) bool {
		return kept[o] || o.Name == "ambience" || o.Name == "bird" || o.Name == "animal"
	}}
	return g, f
}
