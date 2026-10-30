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

// chopJob is felling a tree near the town: standing beside it, facing it,
// striking as many times as the axe takes. Of the trees about, one is
// picked at random, so the woodcutters spread out.
func (m *mind) chopJob(f *farm) *job {
	v := m.v
	c := v.town.Center
	var near []*job
	for r := c.Row - 9; r <= c.Row+9; r++ {
		for col := c.Col - 15; col <= c.Col+15; col++ {
			if f.thingAt(col, r) != tileID["tree0"] {
				continue
			}
			for _, s := range []struct {
				dc, dr int
				face   string
			}{{-1, 0, "right"}, {1, 0, "left"}, {0, -1, "down"}, {0, 1, "up"}} {
				if f.open(col+s.dc, r+s.dr) {
					near = append(near, &job{to: pathfind.Cell{Col: col + s.dc, Row: r + s.dr}, face: s.face, slot: slotAxe, uses: axes[v.kit.axe].hits, rest: 1})
					break
				}
			}
		}
	}
	if len(near) == 0 {
		return nil
	}
	return near[v.luck.Intn(len(near))]
}
