package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/town"
)

// forgeBy is the town whose forge o stands in front of, or nil.
func (f *farm) forgeBy(o *game.Object) *town.Town {
	col, row := f.cellOf(feet(o))
	for _, t := range f.towns.Near(col, row, f.towns.Reach()+4) {
		if row == t.Forge.Row+1 && col >= t.Forge.Col-1 && col <= t.Forge.Col+1 {
			return t
		}
	}
	return nil
}

// bakeBread is the player baking a sheaf of wheat into a loaf at a forge,
// and reports whether they had the wheat.
func (f *farm) bakeBread() bool {
	c, _ := cropIndex("wheat")
	if !f.bag.Take(cropOf(c), 1) {
		return false
	}
	f.bag.Add(bread, 1)
	at := feet(f.player)
	at.Y -= tile
	f.float(at, "+1 "+i18n.T("Bread"), 'y')
	f.sound("clang", &at)
	f.earnBy(f.kit, xpSow)
	return true
}

// forgeNear is the interaction of a forge with wheat in the bag: bake it.
func (f *farm) forgeNear(o *game.Object) (words string, do func(), ok bool) {
	c, _ := cropIndex("wheat")
	if f.forgeBy(o) == nil || f.bag.Count(cropOf(c)) == 0 {
		return "", nil, false
	}
	return i18n.T("E: bake bread"), func() {
		if f.bakeBread() {
			clip, seconds := useClip(actForge)
			f.hands.acting = seconds
			f.player.Animator.Play(f.hands.pose(f.player) + clip)
		}
	}, true
}
