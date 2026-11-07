package farm

import "testing"

// Going to bed saves the game: a new session loading it wakes on the morning
// after, with what the player had.
func TestSleepingSavesTheGame(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	const title = "sge sleep test"
	s := New(SVGArt(), idle{})
	s.EnableSave(title)
	if HasSave(title) {
		t.Fatal("a save before any")
	}
	f := s.f
	f.money = 77
	f.clock.Hour = 21
	day := f.clock.Day
	if !f.sleepAndSave() {
		t.Fatal("did not sleep")
	}
	if !HasSave(title) {
		t.Fatal("no save after sleeping")
	}
	g := New(SVGArt(), idle{})
	if err := g.Load(title); err != nil {
		t.Fatal(err)
	}
	if g.f.money != 77 || g.f.clock.Day != day+1 || g.f.clock.Hour != wakeAt {
		t.Errorf("money %d, day %d, hour %v", g.f.money, g.f.clock.Day, g.f.clock.Hour)
	}
}

// With no saving set up, as in a demo, sleeping still works.
func TestSleepingWithNoSavingOnlySleeps(t *testing.T) {
	_, f := newGame(idle{})
	f.clock.Hour = 21
	if !f.sleepAndSave() || f.clock.Hour != wakeAt {
		t.Errorf("hour %v", f.clock.Hour)
	}
}
