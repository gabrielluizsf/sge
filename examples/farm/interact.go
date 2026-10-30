package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// talkReach is how near a villager one must stand to talk to them.
const talkReach = 1.5 * tile

// interaction is what E does where the player o stands: its words over
// their head, and doing it. ok is false where it does nothing.
func (f *farm) interaction(o *game.Object) (words string, do func(), ok bool) {
	if v := f.nearPerson(o); v != nil {
		return "E: falar com " + v.name, func() { f.talk = &talk{v: v}; f.sound("ui_open", nil) }, true
	}
	if f.room != nil {
		if f.room.own && f.byBed(o) {
			return "E: dormir", func() { f.sleep() }, true
		}
		return "", nil, false
	}
	col, row := f.cellOf(feet(o))
	if d, ok := f.doorAt(col, row); ok {
		return "E: entrar", func() { f.enter(d) }, true
	}
	if s, ok := f.stallBy(o); ok {
		return "E: loja", func() { f.shop.open, f.shop.choice, f.shop.at = true, 0, s; f.sound("ui_open", nil) }, true
	}
	return "", nil, false
}

// nearPerson is the villager nearest the player o within talking reach,
// in the room or out of doors, or nil.
func (f *farm) nearPerson(o *game.Object) *villager {
	var best *villager
	far := talkReach
	consider := func(v *villager, at vector.Vector3) {
		if d := at.Sub(feet(o)).Length(); d <= far {
			best, far = v, d
		}
	}
	if f.room != nil {
		for _, r := range f.residents {
			consider(r.v, r.obj.Transform.Position)
		}
		return best
	}
	for _, v := range f.awake() {
		if !v.obj.Hidden {
			consider(v, v.obj.Transform.Position)
		}
	}
	return best
}

// byBed reports whether o stands beside a bed of the room.
func (f *farm) byBed(o *game.Object) bool {
	for r := range roomRows {
		for c := range roomCols {
			if id := f.room.things.At(c, r); id != tileID["bed"] && id != tileID["bed_blue"] {
				continue
			}
			if f.ground.CellCenter(c, r).Sub(feet(o)).Length() <= talkReach {
				return true
			}
		}
	}
	return false
}
