package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
)

// openVisit takes the player to the plaza of a town at noon on a day when a
// grown villager of the next town stays there, to see a visitor from
// elsewhere.
func (s *Session) openVisit() error {
	f := s.f
	for r := range 6 {
		for c := -3; c < 6; c++ {
			t := f.towns.Of(c, r)
			if t == nil || f.roadFrom(t) == nil {
				continue
			}
			f.settle(t)
			for day := range 40 {
				f.clock.Day = day
				for _, v := range f.people {
					if v.town != t || v.stage() < life.Adult {
						continue
					}
					if sp, ok := f.planOf(v).visit(); ok {
						to := f.roadFrom(t).to
						f.settle(to)
						fmt.Printf("visitor %s (%s) of the town at %v, visiting the one at %v\n", v.name, v.id, t.Center, to.Center)
						f.clock.Hour = (sp.from + sp.to) / 2
						f.hands.ctl = still{}
						f.player.Transform.Position = f.center(pathfind.Cell(to.Center))
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("farm: no one visits another town")
}
