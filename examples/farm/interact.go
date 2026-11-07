package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// talkReach is how near a villager one must stand to talk to them.
const talkReach = 1.5 * tile

// interaction is what E does where the player o stands: its words over
// their head, and doing it. ok is false where it does nothing.
func (f *farm) interaction(o *game.Object) (words string, do func(), ok bool) {
	if f.ride != nil {
		return i18n.T("E: get down"), f.getDown, true
	}
	if f.room == nil {
		if t, name := f.shedBy(o); t != nil {
			return i18n.T(shedWords[name]), func() { f.openShed(t, name) }, true
		}
		if words, do, ok := f.forgeNear(o); ok {
			return words, do, true
		}
	}
	if v := f.nearPerson(o); v != nil {
		return i18n.T("E: talk to %s", v.name), func() { f.openTalk(v) }, true
	}
	if f.room != nil {
		return f.roomInteraction(o)
	}
	if y := f.chestBy(o); y != nil {
		return i18n.T("E: open chest"), func() { f.openChest(&y.p) }, true
	}
	if words, do, ok := f.cowNear(o); ok {
		return words, do, true
	}
	if words, do, ok := f.sheepNear(o); ok {
		return words, do, true
	}
	if words, do, ok := f.horseNear(o); ok {
		return words, do, true
	}
	col, row := f.cellOf(feet(o))
	if d, ok := f.doorAt(col, row); ok {
		return i18n.T("E: enter"), func() { f.enter(d) }, true
	}
	if s, ok := f.stallBy(o); ok {
		return i18n.T("E: shop"), func() {
			f.shop.open, f.shop.kind, f.shop.choice, f.shop.at = true, "", 0, s
			f.sound("ui_open", nil)
		}, true
	}
	if f.mineable(o) {
		return i18n.T("E: mine"), func() { f.hands.useNow = true }, true
	}
	return "", nil, false
}

// mineable reports whether the player o holds the pickaxe with a rock
// where they work.
func (f *farm) mineable(o *game.Object) bool {
	if f.slot != slotPick || f.hands.swims || f.bag.Count(itemPick) == 0 {
		return false
	}
	return isRock(f.thingAt(f.hands.useCell(o)))
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
	return f.nearPiece(o, "bed") >= 0
}

// nearPiece is the index in the room's pieces of the piece called name
// that o stands within reach of, the nearest, or -1.
func (f *farm) nearPiece(o *game.Object, name string) int {
	best, far := -1, talkReach
	for i, p := range f.home {
		if p.name != name {
			continue
		}
		s := p.cells(f.art)
		for r := s.r0; r < s.r0+s.rows; r++ {
			for c := s.c0; c < s.c0+s.cols; c++ {
				if d := f.ground.CellCenter(c, r).Sub(feet(o)).Length(); d <= far {
					best, far = i, d
				}
			}
		}
	}
	return best
}

// roomInteraction is what E does in a room: sleep beside a bed, sit on a
// chair, open a chest, in the player's own house.
func (f *farm) roomInteraction(o *game.Object) (words string, do func(), ok bool) {
	if !f.room.own {
		return "", nil, false
	}
	if f.seat != nil {
		return f.chairNear()
	}
	if f.byBed(o) {
		return i18n.T("E: sleep"), func() { f.goToBed() }, true
	}
	if i := f.nearPiece(o, "chest"); i >= 0 {
		return i18n.T("E: open chest"), func() { f.openChest(&f.home[i]) }, true
	}
	return f.chairNear()
}

// chestBy is the chest out of doors nearest o within reach, or nil.
func (f *farm) chestBy(o *game.Object) *yardChest {
	var best *yardChest
	far := talkReach
	for _, y := range f.yard {
		if d := f.ground.CellCenter(y.p.col, y.p.row).Sub(feet(o)).Length(); d <= far {
			best, far = y, d
		}
	}
	return best
}

// openChest opens the inventory with the chest p beside it.
func (f *farm) openChest(p *piece) {
	f.pack.Side, f.pack.SideTitle = chestStore{f, p}, i18n.T("Chest")
	f.pack.Show()
	f.pack.SetTab(0)
	f.syncPackHint()
	f.sound("ui_open", nil)
}
