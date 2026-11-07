package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// openWork takes the player to a villager at the chore c, at the middle of
// a spell of their plan that has it, to look at them working: the town
// wakes at that hour with them in place, the job in hand.
func (s *Session) openWork(c chore) error {
	f := s.f
	for r := range 6 {
		for col := -3; col < 6; col++ {
			t := f.towns.Of(col, r)
			if t == nil {
				continue
			}
			f.settle(t)
			f.stageRock(t, c)
			for day := range 40 {
				for _, v := range f.people {
					if v.town != t || v.ours || v.stage() < life.Adult {
						continue
					}
					for i, sp := range f.planOn(v, day) {
						if sp.doing != c {
							continue
						}
						if j := f.siteOf(v, c, spellKey(v, day, i)); j != nil {
							f.clock.Day, f.clock.Hour = day, (sp.from+sp.to)/2
							f.hands.ctl = still{}
							f.player.Transform.Position = f.center(pathfind.Cell{Col: j.to.Col - 2, Row: j.to.Row + 1})
							fmt.Printf("working: %s (%s) %+v on day %d at %.1f at %v\n", v.name, v.id, v.who, day, f.clock.Hour, j.to)
							return nil
						}
					}
				}
			}
		}
	}
	return fmt.Errorf("farm: no one at that chore")
}

// stageRock lays a rock by the plaza of town t when the chore is mining, as
// the realistic look keeps the rocks too far from its towns for the demo.
func (f *farm) stageRock(t *town.Town, c chore) {
	if c == choreMine {
		f.setThing(t.Center.Col+3, t.Center.Row-1, tileID[rockName(TreeBiomes[0], 0)])
	}
}
