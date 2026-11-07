package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// Open opens a screen on s, to look at it: "craft", "arrange" or "chest" (see openHouse), "fight" with a villager who strikes back, "sleep" in the house of one of the east town, at night, "visit" in the town a villager of the next one visits, at noon, "inventory" with a bag of a
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
			f.ui.reader.Open(i18n.T("License - %s", c.Name), c.License())
			f.ui.open = licensing
		}
	case "fight":
		if err := s.openFight(d); err != nil {
			return nil, err
		}
	case "sleep", "sleepteen":
		if err := s.openSleep(what == "sleepteen"); err != nil {
			return nil, err
		}
	case "mine", "chop":
		if err := s.openWork(map[string]chore{"mine": choreMine, "chop": choreChop}[what]); err != nil {
			return nil, err
		}
	case "visit":
		if err := s.openVisit(); err != nil {
			return nil, err
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
		if !s.openHouse(what) {
			return nil, fmt.Errorf("farm: no screen called %q", what)
		}
	}
	return s.Play(d)
}
