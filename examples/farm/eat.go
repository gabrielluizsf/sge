package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/vitals"
)

// ripeAt reports whether the plot on cell col, row has a ripe crop, which
// is picked, not eaten.
func (f *farm) ripeAt(col, row int) bool {
	p := f.plotAt(col, row)
	return p != nil && p.ripened()
}

// eat is the player eating one of the crop in hand, if they have one and
// are not full: it feeds them as much as the crop does.
func (f *farm) eat(k *kit) action {
	if k != f.kit {
		return actNone
	}
	at := f.player.Transform.Position
	at.Y -= tile
	if f.vit.Food >= vitals.Max-0.5 {
		f.float(at, i18n.T("Full"), 'r')
		f.sound("deny", nil)
		return actNone
	}
	c := crops[k.seed]
	if !k.bag.Take(cropOf(k.seed), 1) {
		return actNone
	}
	f.vit.Eat(c.nutrition)
	f.float(at, i18n.T("+%d food", int(c.nutrition)), 'g')
	f.sound("harvest", nil)
	return actEat
}
