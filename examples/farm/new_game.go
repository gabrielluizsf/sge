package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/scene"
	"github.com/gabrielluizsf/sge/town"
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
	g.Add(w)
	f.ui = newUI(f)
	menus := game.NewObject("menus")
	menus.Always = true
	menus.AddScript(f.ui)
	g.Add(menus)
	o, p := newFarmer(f, ctl)
	f.player, f.hands = o, p
	o.Transform.Position = standIn(f, where.spawn.Col, where.spawn.Row)
	g.Add(o)
	homeDoor = where.door
	g.Add(f.newStall(farmStall))
	g.Add(f.newHouse(0, town.Cell{Col: where.home.col, Row: where.home.row}, true))
	wild.Props = f.props
	// The wild makes its chunks, and the towns wake, where the camera has
	// just been pointed, so they come after the farmer, whom it follows.
	g.Add(wild.Object)
	f.keeper = &town.Keeper{Layout: towns, Range: 34, Cell: tile, Wake: f.wake}
	k := game.NewObject("keeper")
	k.AddScript(f.keeper)
	g.Add(k)
	night := game.NewObject("night")
	night.Z, night.Screen = zNight, true
	night.AddScript(daytime.Shade{Clock: f.clock,
		Night: canvas.RGBA(canvas.RGBAParams{R: 14, G: 18, B: 52, A: 165}),
		Glow:  canvas.RGBA(canvas.RGBAParams{R: 255, G: 128, B: 60, A: 70})})
	g.Add(night)
	h := game.NewObject("hud")
	h.Z, h.Screen = zHUD, true
	h.AddScript(&hud{f: f, farmer: o, p: p})
	g.Add(h)
	kept := map[*game.Object]bool{w: true, menus: true, o: true, night: true, h: true}
	f.stage = &scene.Stage{Game: g, Current: outside, Keep: func(o *game.Object) bool { return kept[o] || o.Name == "ambience" }}
	return g, f
}
