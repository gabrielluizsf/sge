package farm

import (
	"strconv"

	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// flashFor is how long the screen's edge stays red after a blow, and
// flinchFor how long one struck reels.
const (
	flashFor  = 0.5
	flinchFor = 0.3
)

// wound is villager v struck for harm by someone at from: they bleed and
// reel, and take out an axe to go after the player; with their health run
// out they die, as a villager of old age does.
func (f *farm) wound(v *villager, harm float64, from vector.Vector3) {
	at := feet(v.obj)
	f.float(vector.Vector3{X: at.X, Y: at.Y - tile}, "-"+strconv.Itoa(int(harm)), 'r')
	if v.wounds += harm; v.wounds >= npcHealth {
		f.bleed(at, from, 40)
		f.slay(v)
		return
	}
	f.bleed(at, from, 14)
	v.angry = npcAnger
	if v.kit.axe == noAxe {
		v.kit.axe = 0
	}
	v.kit.slot = slotAxe
	v.mind.drop(0)
	f.reel(v.hands, v.obj, from)
}

// slay is villager v dead of a blow: gone from their town and the game
// like one who died of old age.
func (f *farm) slay(v *villager) {
	f.family.Die(v.id)
	f.leaveLife(v)
	f.notify(i18n.T("%s was killed.", v.name))
	f.sound("deny", nil)
}

// hurtPlayer is the player struck for harm by someone at from: whatever
// they were doing, a talk, a shop, a chair, arranging, ends; they bleed,
// reel and flash red; with their health run out they die.
func (f *farm) hurtPlayer(harm float64, from vector.Vector3) {
	f.talk, f.shop.open = nil, false
	f.stopMilking()
	if f.edit != nil {
		f.stopEdit()
	}
	f.standUp()
	at := feet(f.player)
	f.bleed(at, from, 14)
	f.flash = flashFor
	f.float(vector.Vector3{X: at.X, Y: at.Y - tile}, "-"+strconv.Itoa(int(harm)), 'r')
	f.reel(f.hands, f.player, from)
	if f.vit.Hurt(harm) {
		f.fall()
	}
}

// reel throws p's person o back from the blow that came from from, and
// stops them a moment.
func (f *farm) reel(p *farmer, o *game.Object, from vector.Vector3) {
	if away := feet(o).Sub(from); away.Length() > 0 {
		d := away.Scale(0.35 * tile / away.Length())
		tilemap.Move(o, d.X, d.Y, f.solids()...)
	}
	if p != nil {
		p.acting = flinchFor
	}
}
