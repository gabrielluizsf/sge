package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
)

// chaseEvery is how often, in seconds, a villager going after the player
// finds the way anew, as the player moves.
const chaseEvery = 0.4

// meleeAt is how close a villager must get to the player to swing at them.
const meleeAt = 1.1 * tile

// fight is a villager angry with the player: they walk to them by the way
// pathfind finds, and, near enough, swing their axe at intervals. They give
// up when the player gets far away or goes in a house; the anger fades with
// time.
func (m *mind) fight(f *farm, o *game.Object, dt float64) intent {
	v := m.v
	at := feet(o)
	d := feet(f.player).Sub(at)
	if v.angry -= dt; v.angry <= 0 || f.room != nil || d.Length() > npcLeash {
		v.angry = 0
		m.drop(1)
		return intent{}
	}
	if v.swing -= dt; d.Length() <= meleeAt {
		in := intent{face: faceTo(at, f.player)}
		if v.swing <= 0 {
			v.swing, in.strike = npcEvery, true
		}
		return in
	}
	if m.chase -= dt; m.route == nil || m.chase <= 0 {
		m.chase = chaseEvery
		col, row := f.cellOf(at)
		pc, pr := f.cellOf(feet(f.player))
		if way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, pathfind.Cell{Col: pc, Row: pr}); len(way) > 0 {
			m.route = &pathfind.Route{Cells: way}
		} else {
			m.route = nil
		}
	}
	if m.route != nil {
		if way, done := m.route.Steer(at, f.center); !done {
			return intent{dir: way, upTo: way.Length(), run: true}
		}
	}
	return intent{dir: d, run: true}
}
