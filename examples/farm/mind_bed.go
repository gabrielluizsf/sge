package farm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
)

// bedUnseen puts a villager out of sight to bed when it is their night:
// they are at home, in, as if they had walked there unseen, so whoever
// goes into their house finds them in bed, and not out in the cold.
func (m *mind) bedUnseen(f *farm, o *game.Object) {
	v := m.v
	if v.inside || v.home == nil || v.now(f) != doSleep {
		return
	}
	o.Transform.Position = f.center(pathfind.Cell(v.home.Door))
	o.Hidden, v.inside, v.homeFor = true, true, 0
	m.job, m.route = nil, nil
}
