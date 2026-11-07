package farm

import "testing"

// The film goes through all its shots, a wolf felled, trees struck, milk
// and money made along the way, the farmer never dying, and a new day each
// time round.
func TestTheFilmBehindTheMenu(t *testing.T) {
	g, f := build(SVGArt(), Cinema())
	c := f.hands.ctl.(*cinema)
	day := f.clock.Day
	seen := map[int]bool{}
	wood0, felled, milked, sold := f.bag.Count(wood), false, false, false
	line := f.self.id
	total := 0.0
	for _, s := range film {
		total += s.seconds
	}
	for range int((total + 3) / frame) {
		g.Step(frame)
		seen[c.at] = true
		if f.wolves != nil {
			for _, w := range f.wolves.wolves {
				felled = felled || w.fallen
			}
		}
		milked = milked || f.bag.Count(milk) > 0
		sold = sold || f.kit.sold > 0
		if f.vit.Dead() {
			t.Fatal("the farmer died")
		}
	}
	if len(seen) != len(film) {
		t.Errorf("saw %d of %d shots", len(seen), len(film))
	}
	if f.clock.Day <= day {
		t.Errorf("still day %d", f.clock.Day)
	}
	if !felled {
		t.Error("no wolf was felled")
	}
	if !milked {
		t.Error("no milk")
	}
	if !sold {
		t.Errorf("nothing sold: %d wood from %d", f.bag.Count(wood), wood0)
	}
	if f.self.id != line {
		t.Error("the heir took over")
	}
}

// A cut fades to black and back.
func TestACutFades(t *testing.T) {
	c := &cinema{began: true, time: 0}
	if c.fade() != 1 {
		t.Errorf("fade at the start %v", c.fade())
	}
	c.time = film[0].seconds / 2
	if c.fade() != 0 {
		t.Errorf("fade in the middle %v", c.fade())
	}
	c.time = film[0].seconds
	if c.fade() != 1 {
		t.Errorf("fade at the end %v", c.fade())
	}
}
