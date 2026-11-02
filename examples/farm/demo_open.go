package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// Open opens a screen on s, to look at it: "inventory" with a bag of a
// little of everything, "pause", "credits", "licence" (the first credit's)
// or "talk" with the smith of the east town, the player taken to them; and plays d more of the
// game, so that what is said is typed out. It is an error to ask for any
// other, or for a talk with no one.
func (s *Session) Open(what string, d Demo) (*canvas.Canvas, error) {
	f := s.f
	switch what {
	case "inventory":
		for i := range crops {
			f.bag.Add(seedsOf(i), 3+i)
			f.bag.Add(cropOf(i), 2+i)
		}
		f.bag.Add(wood, 7)
		f.axe = 1
		f.stock(itemAxe)
		f.pack.Show()
	case "pause":
		f.ui.open = pausing
		f.ui.pause.Lines = f.ui.menuLines()
	case "credits", "licence":
		f.ui.open = crediting
		if what == "licence" && len(f.ui.list) > 0 {
			c := f.ui.list[0]
			f.ui.reader.Open("Licença - "+c.Name, c.License())
			f.ui.open = licensing
		}
	case "talk":
		t := f.towns.Of(1, 0)
		f.settle(t)
		var v *villager
		for _, p := range f.people {
			if p.town == t && p.smith() {
				v = p
			}
		}
		if v == nil {
			return nil, fmt.Errorf("farm: no one to talk to")
		}
		f.player.Transform.Position = v.pos
		f.openTalk(v)
	default:
		return nil, fmt.Errorf("farm: no screen called %q", what)
	}
	return s.Play(d)
}
