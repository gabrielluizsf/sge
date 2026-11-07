package farm

import (
	"math"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/town"
)

// settleSpan settles the work of town t's people from hour from to hour to
// of the game's time: each day of it, each villager's spells of it that
// fall within, then the sales between them while they traded together.
func (f *farm) settleSpan(t *town.Town, from, to float64) {
	for day := int(math.Floor(from / 24)); float64(day)*24 < to; day++ {
		var traders []trading
		for _, v := range f.people {
			if v.town != t || v.ours || v.stage() < life.Adult {
				continue
			}
			for i, s := range f.planOn(v, day) {
				lo, hi := max(float64(day)*24+s.from, from), min(float64(day)*24+s.to, to)
				if hi <= lo {
					continue
				}
				f.settleSpell(v, s.doing, spellKey(v, day, i), hi-lo)
				if s.doing == choreTrade {
					traders = append(traders, trading{v, lo, hi})
				}
			}
		}
		f.settleSwaps(traders)
	}
}
