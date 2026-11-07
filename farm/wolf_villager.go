package farm

import (
	"strconv"

	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// wolfSeen is how near a grown villager must be to a wolf to take it on.
const wolfSeen = 9 * tile

// teenAndUp is the youngest a wolf goes for: children are kept in, at dusk,
// by their parents.
const teenAndUp = life.Teen

// bitten is villager v bitten by a wolf at from: they bleed and reel, a
// grown one takes out an axe, and with their health run out they die.
func (f *farm) bitten(v *villager, harm float64, from vector.Vector3) {
	at := feet(v.obj)
	f.float(vector.Vector3{X: at.X, Y: at.Y - tile}, "-"+strconv.Itoa(int(harm)), 'r')
	if v.wounds += harm; v.wounds >= npcHealth {
		f.bleed(at, from, 40)
		f.family.Die(v.id)
		f.leaveLife(v)
		f.notify(i18n.T("%s was killed by a wolf.", v.name))
		return
	}
	f.bleed(at, from, 14)
	f.reel(v.hands, v.obj, from)
}

// wolfNear is the wolf a grown villager at at would take on, or nil.
func (f *farm) wolfNear(v *villager, at vector.Vector3) *wolf {
	if f.wolves == nil || v.stage() < life.Adult {
		return nil
	}
	return f.wolves.nearest(at, wolfSeen)
}

// fightWolf is a villager going for wolf w: out comes the axe, they run to
// it by the way pathfind finds, and swing when it is near enough.
func (m *mind) fightWolf(f *farm, o *game.Object, w *wolf, dt float64) intent {
	v := m.v
	if v.kit.axe == noAxe {
		v.kit.axe = 0
	}
	v.kit.slot = slotAxe
	at := feet(o)
	d := feet(w.obj).Sub(at)
	if v.swing -= dt; d.Length() <= meleeAt {
		in := intent{face: faceTo(at, w.obj)}
		if v.swing <= 0 {
			v.swing, in.strike = npcEvery, true
		}
		return in
	}
	if m.chase -= dt; m.route == nil || m.chase <= 0 {
		m.chase = chaseEvery
		col, row := f.cellOf(at)
		wc, wr := f.cellOf(feet(w.obj))
		if way := f.grid().Path(pathfind.Cell{Col: col, Row: row}, pathfind.Cell{Col: wc, Row: wr}); len(way) > 0 {
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

// killedBy is the player's blow on w: when it fells it, they have the
// credit.
func (f *farm) killedBy(p *wolfPack, w *wolf, harm float64, from vector.Vector3) {
	p.strike(w, harm, from)
	if w.fallen {
		f.earnBy(f.kit, xpWolf)
		f.notify(i18n.T("Wolf killed"))
	}
}
