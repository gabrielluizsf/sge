package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vector"
)

// outside is the scene of the world outside the houses.
const outside = "outside"

// enter takes the player into the house at door: the world outside waits
// as it is, and the room is laid out with whoever of the house is in,
// the player on the mat facing in.
func (f *farm) enter(d houseDoor) {
	r := newRoom(d)
	name := fmt.Sprint("house ", d.at.Col, ",", d.at.Row)
	f.stage.Switch(name, func(g *game.Game) {
		g.Add(r.ground.Object)
		g.Add(r.things.Object)
		g.Covered, g.Background = false, roomDark
		f.residents = nil
		for i, v := range f.people {
			if v.home == nil || v.home.Door != d.at || !v.inside {
				continue
			}
			o := f.resident(v, pathfind.Cell{Col: 2 + i%3, Row: 2 + i%2})
			g.Add(o)
		}
	})
	f.room = r
	f.player.Transform.Position = f.ground.CellCenter(matCol, roomRows-2)
	f.face("up")
	f.sound("ui_open", nil)
}

// leave takes the player out of the room, to stand before its door
// facing away from it, the world outside as it was left.
func (f *farm) leave() {
	r, name := f.room, f.stage.Current
	f.stage.Switch(outside, nil)
	f.stage.Forget(name)
	f.room, f.residents, f.talk = nil, nil, nil
	f.player.Transform.Position = f.ground.CellCenter(r.door.Col, r.door.Row)
	f.face(r.faceOut)
	f.sound("ui_close", nil)
}

// resident is villager v at home, standing on cell c of the room.
func (f *farm) resident(v *villager, c pathfind.Cell) *game.Object {
	o, _ := f.personOf(v)
	o.Transform.Position = f.ground.CellCenter(c.Col, c.Row)
	f.residents = append(f.residents, resident{v: v, obj: o})
	return o
}

// doorAt is the door of a house whose doorstep is cell col, row, the
// player's or a town's.
func (f *farm) doorAt(col, row int) (houseDoor, bool) {
	if col == homeDoor.Col && row == homeDoor.Row {
		return houseDoor{at: homeDoor, faceOut: "down", kind: 0, seed: seed, own: true}, true
	}
	for _, t := range f.towns.Near(col, row, town.Reach+4) {
		for i, h := range t.Houses {
			if h.Door.Col == col && h.Door.Row == row {
				out := "up"
				if h.Row < row {
					out = "down"
				}
				return houseDoor{at: h.Door, faceOut: out, kind: h.Kind, seed: t.Seed + uint32(i)}, true
			}
		}
	}
	return houseDoor{}, false
}

// face turns the player to face dir.
func (f *farm) face(dir string) {
	if d, ok := directions[dir]; ok {
		f.hands.facing = vector.Vector3{X: float64(d[0]), Y: float64(d[1])}
	}
}
