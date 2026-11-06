package farm

import (
	"encoding/json"
	"testing"
)

func TestXPAndLevelRoundTrip(t *testing.T) {
	_, f := newGame(idle{})
	f.earn(xpNeed(1) + xpNeed(2) + 17)
	s := takeSnapshot(f)
	b, _ := json.Marshal(s)
	var back snapshot
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	_, g := newGame(idle{})
	if err := restore(g, &back); err != nil {
		t.Fatal(err)
	}
	if g.prog != f.prog || g.prog.Level != 3 || g.prog.XP != 17 || g.shown != g.xpValue() {
		t.Errorf("loaded %+v shown %v", g.prog, g.shown)
	}
}

func TestAnOldSaveWithoutXPIsLevelOne(t *testing.T) {
	var old map[string]any
	json.Unmarshal(playedJSON(t), &old)
	delete(old, "level")
	delete(old, "xp")
	b, _ := json.Marshal(old)
	var s snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatal(err)
	}
	_, g := newGame(idle{})
	if err := restore(g, &s); err != nil || g.prog.Level != 1 || g.prog.XP != 0 {
		t.Errorf("%v %+v", err, g.prog)
	}
}

func TestHostileXPIsRefused(t *testing.T) {
	for name, edit := range map[string]func(*snapshot){
		"negative xp":      func(s *snapshot) { s.Level, s.XP = 2, -1 },
		"negative level":   func(s *snapshot) { s.Level = -4 },
		"huge level":       func(s *snapshot) { s.Level = 1 << 40 },
		"xp over a level":  func(s *snapshot) { s.Level, s.XP = 1, xpNeed(1) },
		"xp with no level": func(s *snapshot) { s.Level, s.XP = 0, 5 },
		"huge xp":          func(s *snapshot) { s.Level, s.XP = 5, 1<<50 },
	} {
		hostile(t, name, edit)
	}
}
