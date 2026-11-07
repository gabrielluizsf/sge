package farm

import (
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
)

// oddLoad is how much wood a grown villager gathers before selling it
// at the stall.
const oddLoad = 3

// oddJobs is the work of a grown villager with nothing else to do, a
// wanderer's all day: felling trees in sight, and selling what they have
// gathered at the stall. The young have none.
func (m *mind) oddJobs(f *farm) *job {
	v := m.v
	if v.stage() < life.Adult {
		return nil
	}
	if v.kit.bag.Count(wood) >= oddLoad {
		return m.atStall(func(f *farm, at pathfind.Cell) { f.sell(v.kit, f.center(at)) })
	}
	return m.chopJob(f)
}
