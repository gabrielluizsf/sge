package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/pathfind"
)

// The wood a lumberjack keeps for the smith, and how much they gather
// before selling the rest at the stall.
const (
	lumberKeep = 4
	lumberLoad = 10
)

// lumberWork is a lumberjack's work: felling trees about town, and
// selling all but what they keep for the smith once they have a load.
func (m *mind) lumberWork(f *farm) *job {
	k := m.v.kit
	if n := k.bag.Count(wood); n >= lumberLoad {
		return m.atStall(func(f *farm, at pathfind.Cell) {
			sold := k.bag.Count(wood) - lumberKeep
			if sold > 0 && k.bag.Take(wood, sold) {
				k.money += sold * woodPrice
				k.sold += sold
				p := f.center(at)
				f.float(p, fmt.Sprint("+$", sold*woodPrice), 'y')
				f.sound("sell", &p)
			}
		})
	}
	return m.chopJob(f)
}

// chopJob is felling a tree near the town, or, for the pair on duty
// today, around the farm: standing beside it, facing it, striking as many
// times as the axe takes. Of the trees about, one is picked at random, so
// the woodcutters spread out.
func (m *mind) chopJob(f *farm) *job {
	v := m.v
	if f.onDuty(v) {
		around := area{center: pathfind.Cell{Col: f.ground.Cols / 2, Row: f.ground.Rows / 2}, cols: f.ground.Cols/2 + 12, rows: f.ground.Rows/2 + 9}
		if j := m.chopIn(f, around, func(col, row int) bool { return !f.onFarm(col, row) }); j != nil {
			return j
		}
	}
	return m.chopIn(f, area{center: pathfind.Cell(v.town.Center), cols: 15, rows: 9}, func(int, int) bool { return true })
}
