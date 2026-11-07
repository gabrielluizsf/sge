package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/life"
)

// openSleep takes the player into the house of a grown villager of the
// east town at night, whose household is in bed there, to look at them sleep;
// with teen, the child of it is a teen.
func (s *Session) openSleep(teen bool) error {
	f := s.f
	t := f.towns.Of(1, 0)
	f.settle(t)
	for _, v := range f.people {
		if v.town != t || v.home == nil || v.stage() < life.Adult || v.smith() {
			continue
		}
		f.clock.Hour = 23
		for _, m := range f.people {
			if m.town == t && m.home == v.home {
				m.inside = true // the whole household is in bed
				if teen && m.stage() < life.Adult {
					m.who.Stage = life.Teen
				}
				fmt.Printf("asleep: %s %+v\n", m.id, m.who)
			}
		}
		if d, ok := f.doorAt(v.home.Door.Col, v.home.Door.Row); ok {
			f.hands.ctl = still{}
			f.enter(d)
			return nil
		}
	}
	return fmt.Errorf("farm: no one to see asleep")
}
