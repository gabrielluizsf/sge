package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/iso"
)

// newGame is the farm laid out, with the farmer on the path in front of
// the stall, played by ctl.
func newGame(ctl controller) (*game.Game, *farm) {
	g := game.New()
	g.Background = canvas.RGB(43, 31, 46)
	// Seen from a corner: each cell of the world a diamond twice as wide
	// as it is tall.
	g.Projection = iso.New(tile, isoW, tile)
	// The farm and the land around it cover every pixel.
	g.Covered = true
	tiles := newTiles()
	l := newLayers(g, tiles)
	wild := newWild(tiles, l.ground.Cols, l.ground.Rows)
	f := newFarm(l, wild)
	f.ground.AddScript(f)
	o, p := newFarmer(f, ctl)
	o.Transform.Position = standIn(f, stallCol, stallRow+1)
	g.Add(o)
	g.Add(newStall())
	// The wild makes its chunks where the camera has just been pointed,
	// so it comes after the farmer, whom the camera follows.
	g.Add(wild.Object)
	h := game.NewObject("hud")
	h.Z, h.Screen = zHUD, true
	h.AddScript(&hud{f: f, farmer: o, p: p})
	g.Add(h)
	return g, f
}
